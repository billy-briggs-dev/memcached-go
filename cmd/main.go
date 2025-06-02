package main

import (
	"fmt"
	"memcached-go/internal/server"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var port int

	var rootCmd = &cobra.Command{
		Use:   "memcached-go",
		Short: "A simple Golang server using Cobra",
		Run: func(cmd *cobra.Command, args []string) {
			server.Init()
			server.Start(port)
		},
	}

	rootCmd.Flags().IntVarP(&port, "port", "p", 11211, "Port to run the server on")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
