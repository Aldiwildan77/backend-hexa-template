package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = cobra.Command{
	Use:   "backend-hexa-template",
	Short: "A template for building a backend application using Hexagonal Architecture in Go",
	Long:  `This is a template for building a backend application using Hexagonal Architecture in Go. It provides a basic structure and some common functionalities to help you get started quickly.`,
}

var HTTPCmd = cobra.Command{
	Use:   "http",
	Short: "Start the HTTP server",
	Long:  `Start the HTTP server and listen for incoming requests.`,
	Args:  cobra.ArbitraryArgs,
	RunE:  StartHTTPServer,
}

var WorkerCmd = cobra.Command{
	Use:   "worker",
	Short: "Start the worker",
	Long:  `Start the worker and listen for incoming jobs.`,
	Args:  cobra.ArbitraryArgs,
	RunE:  nil,
}

func Execute() {
	cobra.OnInitialize()

	RootCmd.AddCommand(&HTTPCmd)
	RootCmd.AddCommand(&WorkerCmd)

	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
