package result

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"

	"github.com/crewcrew23/proxy-checker/internal/checker"
)

func PrintSummary(results []checker.ProxyResult) {
	alive := 0
	for _, r := range results {
		if r.Alive {
			alive++
			fmt.Printf("✅ %s [%v]\n", r.Proxy, r.Delay)
		} else {
			fmt.Printf("❌ %s [%v]\n", r.Proxy, r.Err)
		}
	}
	fmt.Printf("\nTotal: %d, Alive: %d\n", len(results), alive)
}

func SaveGood(results []checker.ProxyResult, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	wr := csv.NewWriter(file)
	defer wr.Flush()

	sort.Slice(results, func(i, j int) bool {
		return results[i].Delay < results[j].Delay
	})

	for _, r := range results {
		if r.Alive {
			wr.Write([]string{r.Proxy, r.Delay.String()})
		}
	}

	return wr.Error()
}
