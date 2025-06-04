package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/crewcrew23/proxy-checker/internal/checker"
	"github.com/crewcrew23/proxy-checker/internal/loader"
	"github.com/crewcrew23/proxy-checker/internal/result"
)

func main() {
	inputFile := flag.String("input", "", "Path to proxy list file")
	proxyType := flag.String("type", "http", "Proxy type: http, socks5")
	targetDomain := flag.String("target", "http://example.com", "target domain")
	timeout := flag.Int("timeout", 5, "Timeout in seconds")
	saveTo := flag.String("save", "", "Path to save good proxies (optional)")
	flag.Parse()

	if *inputFile == "" {
		log.Fatal("input is required field: example -input input.txt")
	}

	proxies, err := loader.LoadProxies(*inputFile)
	if err != nil {
		log.Fatalf("Failed to load proxies: %v", err)
	}

	fmt.Printf("🔍 Checking %d proxies...\n", len(proxies))
	results := checker.CheckAll(proxies, *targetDomain, *proxyType, *timeout)

	result.PrintSummary(results)

	if *saveTo != "" {
		if err := result.SaveGood(results, *saveTo); err != nil {
			log.Fatalf("Failed to save good proxies: %v", err)
		}
		fmt.Printf("✅ Saved good proxies to %s\n", *saveTo)
	}
}
