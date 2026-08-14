package main

import (
	"fmt"

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
			fmt.Println("Starting the ETA service...")
		},
	}
}