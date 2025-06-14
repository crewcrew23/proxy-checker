package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/crewcrew23/proxy-checker/internal/checker"
	"github.com/crewcrew23/proxy-checker/internal/loader"
	"github.com/crewcrew23/proxy-checker/internal/result"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name: "proxy-checker",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Value:    "",
				Usage:    "file with proxy list",
				Required: true,
			},

			&cli.StringFlag{
				Name:  "type",
				Value: "http",
				Usage: "proxy type http/socks5: default http",
			},

			&cli.StringFlag{
				Name:     "target",
				Value:    "",
				Usage:    "target resource for try access",
				Required: true,
			},

			&cli.IntFlag{
				Name:  "timeout",
				Value: 5,
				Usage: "timeout for request in second: default 5 sec",
			},

			&cli.StringFlag{
				Name:  "save",
				Value: "",
				Usage: "path to save good proxies (options)",
			},
			&cli.StringFlag{
				Name:  "savetype",
				Value: "json",
				Usage: "file extension json/csv (default: json)",
			},
			&cli.IntFlag{
				Name:  "threshold",
				Value: 100,
				Usage: "threshold of the number of proxies in the list, upon reaching which the worker pool will be used for processing (default 100)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {

			input := cmd.String("input")
			proxyType := cmd.String("type")
			target := cmd.String("target")
			timeout := cmd.Int("timeout")
			save := cmd.String("save")
			saveType := cmd.String("savetype")
			threshold := cmd.Int("threshold")

			validProxyType := map[string]bool{"http": true, "socks5": true}
			if !validProxyType[proxyType] {
				return fmt.Errorf("unsupported proxy type: %s", proxyType)
			}

			validFileType := map[string]bool{"json": true, "csv": true}
			if !validFileType[saveType] {
				return fmt.Errorf("unsupported file type: %s", saveType)
			}

			proxies, err := loader.LoadProxies(input)
			if err != nil {
				log.Fatalf("Failed to load proxies: %v", err)
			}

			fmt.Printf("🔍 Checking %d proxies...\n", len(proxies))
			results := checker.CheckAll(proxies, target, proxyType, timeout, threshold)
			result.PrintSummary(results)

			if save != "" {
				if err := result.SaveGood(results, save, saveType); err != nil {
					log.Fatalf("Failed to save good proxies: %v", err)
				}
				fmt.Printf("✅ Saved good proxies to %s\n", save)
			}

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}
