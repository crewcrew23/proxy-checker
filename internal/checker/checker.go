package checker

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

type ProxyResult struct {
	Proxy string
	Alive bool
	Delay time.Duration
	Err   error
}

func CheckAll(proxies []string, target string, proxyType string, timeoutSec int) []ProxyResult {
	var wg sync.WaitGroup
	ch := make(chan ProxyResult, len(proxies))

	for _, p := range proxies {
		wg.Add(1)
		go func(proxyAddr string) {
			defer wg.Done()
			resutl := checkOne(proxyAddr, proxyType, target, timeoutSec)
			ch <- resutl
		}(p)
	}
	wg.Wait()
	close(ch)

	var results []ProxyResult
	for p := range ch {
		results = append(results, p)
	}
	return results
}

func checkOne(proxyAddr, proxyType, target string, timeoutSec int) ProxyResult {
	timeout := time.Duration(timeoutSec) * time.Second
	start := time.Now()

	var client *http.Client

	//input host:port or host:port:login:password
	host, port, auth, err := parseProxyAddr(proxyAddr)
	if err != nil {
		return ProxyResult{Proxy: proxyAddr, Alive: false, Err: err}
	}
	address := net.JoinHostPort(host, port)

	switch proxyType {
	case "http":
		proxyURL := &url.URL{
			Scheme: "http",
			Host:   address,
		}

		if auth != nil {
			proxyURL.User = url.UserPassword(auth.User, auth.Password)
		}

		client = &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSHandshakeTimeout: 10 * time.Second,
				Proxy:               http.ProxyURL(proxyURL),
			},
		}

		return sendHttp(client, proxyAddr, target, start)

	case "socks5":
		dialer, err := proxy.SOCKS5("tcp", address, auth, proxy.Direct)
		if err != nil {
			return ProxyResult{Proxy: proxyAddr, Alive: false, Err: err}
		}

		dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}

		transport := &http.Transport{DialContext: dialContext}
		client = &http.Client{
			Timeout:   timeout,
			Transport: transport,
		}

		return sendSocks5(client, proxyAddr, target, start)

	default:
		return ProxyResult{Proxy: proxyAddr, Alive: false, Err: fmt.Errorf("unsupported proxy type")}
	}

}

func parseProxyAddr(proxyAddr string) (host, port string, auth *proxy.Auth, err error) {
	parts := strings.Split(proxyAddr, ":")
	switch len(parts) {
	case 2:
		host, port = parts[0], parts[1]
	case 4:
		host, port = parts[0], parts[1]
		auth = &proxy.Auth{User: parts[2], Password: parts[3]}
	default:
		err = fmt.Errorf("invalid proxy format: %s", proxyAddr)
	}
	return
}
