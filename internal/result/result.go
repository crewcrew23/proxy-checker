package result

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"

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

func SaveGood(results []checker.ProxyResult, path, format string) error {
	file, err := os.Create(path + "." + format)
	if err != nil {
		return err
	}
	defer file.Close()

	var goodResult []checker.ProxyResult
	for _, p := range results {
		if p.Alive {
			goodResult = append(goodResult, p)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Delay < results[j].Delay
	})

	switch format {
	case "csv":
		wr := csv.NewWriter(file)
		defer wr.Flush()

		for _, r := range results {
			if r.Alive {
				delayToString := strconv.Itoa(int(r.Delay))
				wr.Write([]string{r.Proxy, delayToString})
			}
		}

		return wr.Error()

	case "json":
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", " ")
		return encoder.Encode(goodResult)

	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

}
