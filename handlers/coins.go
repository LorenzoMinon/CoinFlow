package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/LorenzoMinon/coinflow/db"
	"github.com/LorenzoMinon/coinflow/pipeline"
	"github.com/jackc/pgx/v5"
)

type PipelineResult struct {
	Processed int `json:"proccesed"`
	Failed    int `json:"failed"`
}

type Handler struct {
	DB *pgx.Conn
}
type CoinResponse struct {
	ID        int     `json:"id"`
	CoinID    string  `json:"coin_id"`
	Symbol    string  `json:"symbol"`
	PriceUSD  float64 `json:"price_usd"`
	Change24h float64 `json:"change_24h"`
	Sentiment string  `json:"sentiment"`
	FetchedAt string  `json:"fetched_at"`
}

func (h *Handler) RunPipeline(w http.ResponseWriter, r *http.Request) {
	coins := []string{"bitcoin", "ethereum"}
	fetchResults := pipeline.RunWorkerPool(coins, 2)

	processed := 0
	failed := 0

	for _, result := range fetchResults {
		if result.Err != nil {
			failed++
			continue
		}
		record := pipeline.Transform(result)
		err := db.InsertCoinRecord(h.DB, record)
		if err != nil {
			failed++
			continue
		}
		processed++
	}
	summary := PipelineResult{Processed: processed, Failed: failed}
	data, err := json.Marshal(summary)
	if err != nil {
		http.Error(w, "error marshing", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *Handler) GetCoins(w http.ResponseWriter, r *http.Request) {
	my_query := "SELECT * FROM coin_prices ORDER BY fetched_at DESC"
	rows, err := h.DB.Query(context.Background(), my_query)
	if err != nil {
		http.Error(w, "error querying", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var coins []CoinResponse

	for rows.Next() {
		var c CoinResponse
		rows.Scan(&c.ID, &c.CoinID, &c.Symbol, &c.PriceUSD, &c.Change24h, &c.Sentiment, &c.FetchedAt)
		coins = append(coins, c)
	}
	data, err := json.Marshal(coins)
	if err != nil {
		http.Error(w, "error parsing to json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}

func (h *Handler) GetLatest(w http.ResponseWriter, r *http.Request) {
	my_query := "SELECT DISTINCT ON (coin_id) coin_id, symbol, price_usd, change_24h, sentiment, fetched_at FROM coin_prices ORDER BY coin_id, fetched_at DESC"
	rows, err := h.DB.Query(context.Background(), my_query)
	if err != nil {
		http.Error(w, "error querying", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var coins []CoinResponse

	for rows.Next() {
		var c CoinResponse
		rows.Scan(&c.CoinID, &c.Symbol, &c.PriceUSD, &c.Change24h, &c.Sentiment, &c.FetchedAt)
		coins = append(coins, c)
	}
	data, err := json.Marshal(coins)
	if err != nil {
		http.Error(w, "error parsing to json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}
