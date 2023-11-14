package main

import (
	"context"
	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/log/level"
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
	// configure logger
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"svc", "user",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}
	// initialize config
	cfg, err := config.InitConfig(".env")
	if err != nil {
		level.Error(logger).Log(ctx, err, "cannot init config")
		os.Exit(1)
	}

	rootCmd := newRootCmd()
	rootCmd.AddCommand(
		newGRPCServerCmd(cfg, logger),
		newHttpServerCommand(cfg, logger),
	)

	if err = rootCmd.Execute(); err != nil {
		level.Error(logger).Log(ctx, err, "Error executing root command")
		os.Exit(1)
	}

}
