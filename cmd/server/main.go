package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/steffnova/promotions-storage/pkg/api"
	"github.com/steffnova/promotions-storage/pkg/csv"
	"github.com/steffnova/promotions-storage/pkg/storage"

	"github.com/gorilla/mux"
)

func main() {
	flagSet := flag.NewFlagSet("server", flag.ExitOnError)
	flagPeriod := flagSet.Duration("period", time.Second, "period between promotion storage updates")
	flagEnableLog := flagSet.Bool("enable-log", false, "enables/disables logging")
	flagFilePath := flagSet.String("file-path", "promotions.csv", "specify path to csv file from which promotions will be loaded")
	flagSet.Parse(os.Args[1:])

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	storage := storage.WithOptions(
		storage.InMemory(),
		storage.OptionLogging(os.Stdout, *flagEnableLog),
		storage.OptionPeriodicUpdate(ctx, *flagPeriod),
	)

	if err := storage.LoadPromotions(csv.PromotionStreamer(csv.FileReader(*flagFilePath))); err != nil {
		fmt.Printf("Failed to load storage data: %s", err)
		return
	}

	multiplexer := mux.NewRouter()
	multiplexer.HandleFunc("/promotions/{id}", api.PromotionGET(storage.GetPromotion, json.Marshal)).Methods(http.MethodGet)

	server := http.Server{
		Addr:    ":8080",
		Handler: multiplexer,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error while serving http request: %s\n", err)
			done()
		}
	}()

	<-ctx.Done()
	err := server.Shutdown(context.Background())
	switch err {
	case http.ErrServerClosed, nil:
		fmt.Println("Server closed")
	default:
		fmt.Printf("Server closed. Error when closing server: %s\n", err)
	}
}
