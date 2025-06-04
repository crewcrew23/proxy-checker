package checker

import (
	"encoding/base64"
	"net/http"
	"time"
)

func basicAuthHeader(user, pass string) string {
	auth := user + ":" + pass
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func sendSocks5(client *http.Client, proxyAddr, target string, start time.Time) ProxyResult {
	resp, err := client.Get(target)
	if err != nil {
		return ProxyResult{Proxy: proxyAddr, Alive: false, Err: err}
	}
	defer resp.Body.Close()

	delay := time.Since(start)
	return ProxyResult{Proxy: proxyAddr, Alive: true, Delay: delay}
}

func sendHttp(client *http.Client, proxyAddr, target string, start time.Time) ProxyResult {
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return ProxyResult{Proxy: proxyAddr, Alive: false, Err: err}
	}

	resp, err := client.Do(req)
	if err != nil {
		return ProxyResult{Proxy: proxyAddr, Alive: false, Err: err}
	}
	defer resp.Body.Close()

	delay := time.Since(start)
	return ProxyResult{Proxy: proxyAddr, Alive: true, Delay: delay}
}
