package main

import (
	"DaysWithoutFoults/controller"
	counterlib "DaysWithoutFoults/counter"
	"context"
	"errors"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	counter := counterlib.New()

	ctrl := controller.NewController(&counter)

	ticker := time.NewTicker(24 * time.Hour)

	router := http.NewServeMux()
	router.HandleFunc("/", ctrl.GetDaysCount)
	router.HandleFunc("/fail", ctrl.Reset)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	port := os.Getenv("SERVER_PORT")
	server := &http.Server{
		Addr:    ":" + port,
		Handler: c.Handler(router),
	}

	go func() {
		log.Printf("server listening at %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	for {
		select {
		case <-ticker.C:
			counter.Inc()
		case <-ctx.Done():
			log.Printf("exiting...")
			ticker.Stop()
			return
		}
	}
}