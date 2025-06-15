package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/d5avard/factory/internal"
	"github.com/d5avard/factory/internal/requester"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "requester",
	Short: "Run the requester HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	cobra.OnInitialize(func() {
		viper.AddConfigPath("./config")
		viper.SetConfigName("requester")
		viper.SetConfigType("toml")
		if err := viper.ReadInConfig(); err == nil {
			log.Println("Using config file:", viper.ConfigFileUsed())
		}
	})

	rootCmd.PersistentFlags().Int("port", 8000, "Server port")
	rootCmd.PersistentFlags().String("host", "localhost", "Port to run the server on")

	_ = viper.BindPFlag("server.port", rootCmd.PersistentFlags().Lookup("port"))
	_ = viper.BindPFlag("server.host", rootCmd.PersistentFlags().Lookup("host"))

	viper.SetEnvPrefix("REQUESTER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
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

	type Config struct {
		Server struct {
			Port int
		}
		Chatgtp struct {
			APIKey string `mapstructure:"api_key"`
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("unmarshal error: %v", err)
	}

	port := viper.GetInt("server.port")
	host := viper.GetString("server.host")

	router := mux.NewRouter()
	routes := []internal.RouteInjector{
		requester.StatusRoutes{},
		requester.MessagesRoutes{},
	}

	server := internal.NewServer(logger, router, routes)
	addr := fmt.Sprintf("%s:%d", host, port)

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
