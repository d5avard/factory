package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/d5avard/factory/internal"
	"github.com/d5avard/factory/internal/requester"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var port string

var rootCmd = &cobra.Command{
	Use:   "requester",
	Short: "Run the requester HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.Flags().StringVar(&port, "port", "80", "Port to run the server on")
}

type Response struct {
	Message string `json:"message"`
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run() {
	logger := internal.NewLogger("logs", "requester.log")
	defer logger.Close()

	// Example logs
	// logger.Info("Service started", zap.String("env", "production"), zap.Time("started_at", time.Now()))
	// logger.Warn("Disk space low", zap.Int("available_MB", 500))
	// logger.Error("Failed to connect to DB", zap.String("host", "localhost"), zap.Int("port", 5432))

	router := mux.NewRouter()
	routes := []internal.RouteInjector{
		requester.StatusRoutes{},
		requester.MessagesRoutes{},
	}

	server := internal.NewServer(logger, router, routes)
	addr := fmt.Sprintf("localhost:%s", port)

	go func() {
		if err := server.Start(addr); err != nil {
			logger.Error("Server failed to start", zap.String("address", addr), zap.Error(err))
			os.Exit(1)
		}
	}()

	// Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Forced shutdown", zap.Error(err))
	} else {
		logger.Info("Server shutdown gracefully")
	}
}
