package main

import (
	"context"
	logger "user-management/internal/util"

	"github.com/spf13/cobra"
	"os"
	"user-management/internal/config"
)

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "user-management",
		Short: "The legendary user-management service",
	}
}

func main() {

	ctx := context.Background()
	lgr := logger.CreateLogger("user")

	// configure logger

	// initialize config
	cfg, err := config.InitConfig(".env")
	if err != nil {
		logger.Error(lgr, "", ctx, err, "cannot init config")
		os.Exit(1)
	}

	rootCmd := newRootCmd()
	rootCmd.AddCommand(
		newGRPCServerCmd(cfg, lgr),
		newHttpServerCommand(cfg, lgr),
	)

	if err = rootCmd.Execute(); err != nil {
		logger.Error(lgr, "", ctx, err, "Error executing root command")
		os.Exit(1)
	}

}
