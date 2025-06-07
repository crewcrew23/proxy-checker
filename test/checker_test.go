package test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/crewcrew23/proxy-checker/internal/checker"
	"github.com/crewcrew23/proxy-checker/test/mock"
)

func prepareHttp() (*httptest.Server, *httptest.Server) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	proxyMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Proxy-Authorization")
		if auth != "" && auth != "Basic dXNlcjpwYXNz" {
			w.WriteHeader(http.StatusProxyAuthRequired)
			return
		}

		client := http.Client{}
		req, _ := http.NewRequest("GET", ts.URL, nil)
		res, err := client.Do(req)

		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		w.WriteHeader(res.StatusCode)
	}))

	return ts, proxyMock
}

func prepareSocks5() (ts *httptest.Server, addr1 string, close1 func(), addr2 string, close2 func()) {
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	addr1, close1 = mock.StartMockSocks5Server(false)
	addr2, close2 = mock.StartMockSocks5Server(true)

	return
}

func TestChecker_HTTP(t *testing.T) {

	ts, proxyMock := prepareHttp()
	defer ts.Close()
	defer proxyMock.Close()

	proxyURL, _ := url.Parse(proxyMock.URL)
	proxyAddr := proxyURL.Host

	tests := []struct {
		Name      string
		Proxy     string
		ProxyType string
		Expected  bool
	}{
		{
			Name:      "Valid HTTP",
			Proxy:     proxyAddr,
			ProxyType: "http",
			Expected:  true,
		},
		{
			Name:      "Invalid HTTP",
			Proxy:     "invalid:8888",
			ProxyType: "http",
			Expected:  false,
		},
		{
			Name:      "HTTP with valid auth",
			Proxy:     "user:pass@" + proxyAddr,
			ProxyType: "http",
			Expected:  true,
		},
		{
			Name:      "HTTP with invalid auth",
			Proxy:     "wrong:creds@" + proxyAddr,
			ProxyType: "http",
			Expected:  false,
		},
		{
			Name:      "Timeout proxy",
			Proxy:     proxyURL.Host + ":9999",
			ProxyType: "http",
			Expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			result := checker.CheckOne(tt.Proxy, tt.ProxyType, ts.URL, 5)
			if result.Alive != tt.Expected {
				t.Errorf("%s got Alive = %v but Expected %v", tt.Name, result.Alive, tt.Expected)
			}
		})
	}
}

func TestChecker_SOCKS5(t *testing.T) {

	ts, addr1, close1, addr2, close2 := prepareSocks5()
	defer close1()
	defer close2()

	proxyHost := strings.Split(addr1, ":")[0]

	tests := []struct {
		Name      string
		Proxy     string
		ProxyType string
		Expected  bool
	}{
		{
			Name:      "Valid SOCKS5",
			Proxy:     addr1,
			ProxyType: "socks5",
			Expected:  true,
		},
		{
			Name:      "Invalid SOCKS5",
			Proxy:     "invalid:8888",
			ProxyType: "socks5",
			Expected:  false,
		},
		{
			Name:      "SOCKS5 with valid auth",
			Proxy:     "user:pass@" + addr2,
			ProxyType: "socks5",
			Expected:  true,
		},
		{
			Name:      "SOCKS5 with invalid auth",
			Proxy:     "wrong:creds@" + addr2,
			ProxyType: "socks5",
			Expected:  false,
		},
		{
			Name:      "Timeout proxy",
			Proxy:     proxyHost + ":9999",
			ProxyType: "socks5",
			Expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			result := checker.CheckOne(tt.Proxy, tt.ProxyType, ts.URL, 5)
			if result.Alive != tt.Expected {
				t.Errorf("%s got Alive = %v but Expected %v", tt.Name, result.Alive, tt.Expected)
			}
		})
	}
}
