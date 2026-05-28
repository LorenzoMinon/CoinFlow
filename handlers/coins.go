package handlers

import (
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
