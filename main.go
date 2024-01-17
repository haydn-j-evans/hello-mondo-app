package main

import (
	"context"
	"errors"
	"log"
	http "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	log := &log.Logger{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from Mondoo Engineer!"))
	})

	srv := &http.Server{Addr: ":" + "8080"}

	go cleanShutdownListener(srv, log)

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server error: %v", err)
	}

}

func cleanShutdownListener(srv *http.Server, log *log.Logger) {
	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-sigChan
	log.Printf("Received signal %s, shutting down gracefully...", sig)

	// Create a context to attempt a graceful shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt the graceful shutdown by closing the listener and
	// completing all inflight requests.
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Unable to gracefully shutdown server: %v", err)
	}
}
