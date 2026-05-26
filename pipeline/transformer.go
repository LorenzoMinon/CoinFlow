package pipeline

type CoinRecord struct {
	CoinID    string
	Symbol    string
	PriceUSD  float64
	Change24h float64
	Sentiment string
}

func Transform(result FetchResult) CoinRecord {
	sentiment := ""
	if result.Data.Change24h > 2.0 {
		sentiment = "bullish"
	} else if result.Data.Change24h < -2.0 {
		sentiment = "bearish"
	} else {
		sentiment = "neutral"
	}
	return CoinRecord{
		CoinID:    result.CoinID,
		Symbol:    result.CoinID,
		PriceUSD:  result.Data.Price,
		Change24h: result.Data.Change24h,
		Sentiment: sentiment,
	}
}
