package scheduler

import (
	"fmt"
	"time"

	"github.com/LorenzoMinon/coinflow/db"
	"github.com/LorenzoMinon/coinflow/pipeline"
	"github.com/jackc/pgx/v5"
)

func Start(conn *pgx.Conn, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			coins := []string{"bitcoin", "ethereum"}
			fetchResults := pipeline.RunWorkerPool(coins, 5)
			for _, result := range fetchResults {
				if result.Err != nil {
					continue
				}
				record := pipeline.Transform(result)
				db.InsertCoinRecord(conn, record)
				fmt.Println("scheduler running succesfully")
			}
		}
	}()
}
