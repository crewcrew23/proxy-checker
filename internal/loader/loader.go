package loader

import (
	"bufio"
	"os"
	"strings"
)

func LoadProxies(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			proxies = append(proxies, line)
		}
	}

	return proxies, scanner.Err()
}
