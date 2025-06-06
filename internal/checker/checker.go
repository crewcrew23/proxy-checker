package checker

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

type ProxyResult struct {
	Proxy string
	Alive bool
	Delay time.Duration
	Err   error
}

func CheckAll(proxies []string, target, proxyType string, timeoutSec, threshold int) []ProxyResult {
	var outputCH chan ProxyResult

	if len(proxies) > threshold {
		workers := 4 * runtime.NumCPU()
		outputCH = StartCheckerPool(proxies, target, proxyType, timeoutSec, workers)
	} else {
		outputCH = StartSyncGorutines(proxies, target, proxyType, timeoutSec)
	}

	var results []ProxyResult
	for p := range outputCH {
		results = append(results, p)
	}

	return results
}

func CheckOne(proxyAddr, proxyType, target string, timeoutSec int) ProxyResult {
	timeout := time.Duration(timeoutSec) * time.Second
	start := time.Now()

	var client *http.Client

	u, err := ParseProxyString(proxyAddr, proxyType)
	if err != nil {
		return ProxyResult{Proxy: proxyAddr, Alive: false, Err: err}
	}

	address := u.Host
	var auth *proxy.Auth

	if u.User != nil {
		password, _ := u.User.Password()
		auth = &proxy.Auth{
			User:     u.User.Username(),
			Password: password,
		}
	}

	switch proxyType {
	case "http":
		client = &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				TLSHandshakeTimeout: 10 * time.Second,
				Proxy:               http.ProxyURL(u),
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

func ParseProxyString(s string, defaultProtocol string) (*url.URL, error) {
	if !strings.Contains(s, "://") {
		s = defaultProtocol + "://" + s
	}
	return url.Parse(s)
}
