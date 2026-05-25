package pipeline

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CoinData struct {
	Price     float64 `json:"usd"`
	Change24h float64 `json:"usd_24h_change"`
}

type FetchResult struct {
	CoinID string
	Data   CoinData
	Err    error // sends err within the channel instead of losing it
}

func FetchCoin(CoinID string, ch chan FetchResult) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd&include_24hr_change=true", CoinID)
	response, err := http.Get(url)
	if err != nil {
		ch <- FetchResult{CoinID: CoinID, Err: err}
		return
	}
	defer response.Body.Close()

	content, _ := io.ReadAll(response.Body)
	var result map[string]CoinData // coingecko returns a map
	err = json.Unmarshal(content, &result)
	if err != nil {
		ch <- FetchResult{CoinID: CoinID, Err: err}
		return
	}
	coin := result[CoinID]
	ch <- FetchResult{CoinID: CoinID, Data: coin, Err: nil}

}
