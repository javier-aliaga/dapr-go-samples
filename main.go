package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/javier-aliaga/dapr-go-samples/api"
	"github.com/javier-aliaga/dapr-go-samples/dapr"
	"github.com/javier-aliaga/dapr-go-samples/telemetry"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Dapr Workflow runtime (separate goroutine)
	workflowRuntime, err := dapr.StartWorkflowRuntime(ctx)
	if err != nil {
		log.Fatalf("failed to start workflow runtime: %v", err)
	}

	shutdownFn, err := telemetry.Init(ctx, "workflow-app")
	if err != nil {
		log.Fatalf("failed to initialize telemetry: %v", err)
	}

	defer func() {
		_ = shutdownFn(ctx)
	}()

	// Setup HTTP routes
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, workflowRuntime)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Println("HTTP server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	// Wait for signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Println("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}