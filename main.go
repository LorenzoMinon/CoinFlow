package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/LorenzoMinon/coinflow/db"
	"github.com/LorenzoMinon/coinflow/handlers"
	"github.com/LorenzoMinon/coinflow/middleware"
	"github.com/LorenzoMinon/coinflow/scheduler"
)

func main() {
	conn, err := db.Connect()
	if err != nil {
		fmt.Println("Failed db init conn", err)
		os.Exit(1)
	}
	h := handlers.Handler{DB: conn}
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("GET /coins", h.GetCoins)
	http.HandleFunc("GET /coins/latest", h.GetLatest)
	http.HandleFunc("POST /pipeline/run", h.RunPipeline)
	scheduler.Start(conn, 1*time.Minute)
	http.ListenAndServe(":8080", middleware.Logger(http.DefaultServeMux))

}
