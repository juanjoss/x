package main

import (
	"log"
	"os"
	"time"

	"github.com/juanjoss/x/pkg/api"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// setting default flags
	flagset := pflag.NewFlagSet("xflags", pflag.ContinueOnError)
	flagset.Int("port", 8080, "Server port")
	flagset.Duration("http-server-timeout", 30*time.Second, "Server read and write timeout")
	flagset.Duration("server-shutdown-timeout", 5*time.Second, "Server graceful shutdown timeout")

	// init logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// parsing flags
	err := flagset.Parse(os.Args[1:])
	if err != nil {
		log.Printf("error parsing flags: %v", err)
		flagset.PrintDefaults()
		os.Exit(2)
	}

	// binding flags to server config
	viper.BindPFlags(flagset)
	hostname, _ := os.Hostname()
	viper.Set("hostname", hostname)
	viper.AutomaticEnv()

	var config api.Config
	if err := viper.Unmarshal(&config); err != nil {
		logger.Fatal("unable to unmarshal config", zap.Error(err))
	}

	// running the HTTP server
	server := api.NewServer(&config, logger)
	server.ListenAndServe()
}
