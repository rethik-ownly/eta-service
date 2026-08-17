package main

import (
	"github.com/nutanalabs/eta-service/internal/config"
	etaservice "github.com/nutanalabs/eta-service/internal/eta-service"
	"github.com/nutanalabs/eta-service/internal/eta-service/repository"
	"github.com/nutanalabs/eta-service/internal/eta-service/service"
	"github.com/nutanalabs/eta-service/internal/server"
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
			cfg := config.InitConfig("test")
			
			logger.Init(cfg.Log.Level)

			// TODO : dependency injection using wire
			
			repo := repository.NewRepository()
			service := service.NewService(repo)
			handler := etaservice.NewHandler(service)

			handlers := server.Handlers {
				ETAHandler: handler,
			}

			srv := server.NewServer(cfg)

			srv.Run(handlers)
		},
	}
}