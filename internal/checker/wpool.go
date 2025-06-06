package checker

import "sync"

func checkerPool(id int, input <-chan string, output chan<- ProxyResult, target string, proxyType string, timeoutSec int, wg *sync.WaitGroup) {
	defer wg.Done()

	for proxyAddr := range input {
		result := CheckOne(proxyAddr, proxyType, target, timeoutSec)
		output <- result
	}
}

func StartCheckerPool(proxies []string, target string, proxyType string, timeoutSec int, countOfWorker int) chan ProxyResult {
	inputCH := make(chan string, countOfWorker)
	outputCH := make(chan ProxyResult, len(proxies))
	var wg sync.WaitGroup

	for i := 0; i < countOfWorker; i++ {
		wg.Add(1)
		go checkerPool(i, inputCH, outputCH, target, proxyType, timeoutSec, &wg)
	}

	go func() {
		for _, val := range proxies {
			inputCH <- val
		}
		close(inputCH)
	}()

	wg.Wait()
	close(outputCH)

	return outputCH

}

func StartSyncGorutines(proxies []string, target string, proxyType string, timeoutSec int) chan ProxyResult {
	var wg sync.WaitGroup
	outputCH := make(chan ProxyResult, len(proxies))

	for _, p := range proxies {
		wg.Add(1)
		go func(proxyAddr string) {
			defer wg.Done()
			resutl := CheckOne(proxyAddr, proxyType, target, timeoutSec)
			outputCH <- resutl
		}(p)
	}
	wg.Wait()
	close(outputCH)

	return outputCH
}
