package pipeline

func RunWorkerPool(coins []string, numWorkers int) []FetchResult {
	results := make(chan FetchResult, len(coins))
	//job channel
	jobs := make(chan string, len(coins))
	// workers, go routines that read from the chan and process

	for w := 0; w < numWorkers; w++ {
		go func() {
			for CoinID := range jobs {
				FetchCoin(CoinID, results) // where workers send results
			}
		}()
	}
	for _, coin := range coins {
		jobs <- coin
	}
	close(jobs)
	var fetchResults []FetchResult
	for i := 0; i < len(coins); i++ {
		fetchResults = append(fetchResults, <-results)
	}
	return fetchResults
}
