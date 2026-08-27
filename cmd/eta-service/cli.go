package main

import (
	"github.com/nutanalabs/eta-service/internal/config"
	// etaservice "github.com/nutanalabs/eta-service/internal/eta-service"
	// "github.com/nutanalabs/eta-service/internal/eta-service/repository"
	// "github.com/nutanalabs/eta-service/internal/eta-service/service"
	// "github.com/nutanalabs/eta-service/internal/server"
	logger "github.com/nutanalabs/rapido-logger-go"
	"github.com/spf13/cobra"
)


func initCLI() *cobra.Command {
	rootCmd := &cobra.Command {
		Use: "eta-service",
		Short: "A CLI for the ETA service",
	}

	rootCmd.AddCommand(startCommand())

	return rootCmd
}

func startCommand() *cobra.Command {
	return &cobra.Command {
		Use: "start",
		Short: "Starts the ETA service",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize the config
			configFile := "application"
			if len(args) > 0 {
				configFile = args[0]
			}
			cfg := config.InitConfig(configFile)
			logger.Init(cfg.Log.Level)

			// Dependecy injection
			serverDependecies , _ := InitDependencies()

			serverDependecies.server.Run(serverDependecies.handlers)
		},
	}
}