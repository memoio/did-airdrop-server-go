package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/did-server/internal/server"
	"github.com/spf13/cobra"
)

var (
	port  string
	chain string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the application",
	Run: func(cmd *cobra.Command, args []string) {
		srv := server.NewServer(chain, port)

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		fmt.Println("Shutting down server...")
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Server commands",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	runCmd.Flags().StringVarP(&port, "port", "p", "8080", "listen port")
	runCmd.Flags().StringVarP(&chain, "chain", "c", "dev", "listen port")
	ServerCmd.AddCommand(runCmd, stopCmd)
}
