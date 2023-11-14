package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	user "user-management/internal"
	"user-management/internal/middleware"

	"golang.org/x/sync/errgroup"
	"user-management/internal/config"
	"user-management/internal/config/db/postgress"
	"user-management/internal/dao/repo"
	userServ "user-management/internal/implementation"
	usertansport "user-management/internal/transport"
	httpusertransport "user-management/internal/transport/http"
	logger "user-management/internal/util"
)

const (
	shutdownTimeout = time.Second * 5
	readTimeout     = time.Second
	writeTimeout    = time.Second * 10
)

// newHttpServerCommand creates a new cobra command for starting the HTTP server
func newHttpServerCommand(cfg config.Config, lgr log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "http-server",
		Short: "Starts the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {

			log.With(lgr, "newHttpServerCommand")

			// Create a context that is canceled when Interrupt, SIGINT or SIGTERM is received
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
			defer cancel()

			// DB creation
			var db *gorm.DB
			{
				var err error
				// Connect to the "ordersdb" database
				db, err = postgress.InitDB(cfg.DB, lgr)
				if err != nil {
					logger.Error(lgr, "exit", err)
					os.Exit(-1)
				}
			}

			// Create Order Service
			var svc user.Service
			{
				repository, err := repo.NewUserSQLRepository(db, lgr)
				if err != nil {
					logger.Error(lgr, "exit", err)
					os.Exit(-1)
				}
				svc, err = userServ.NewService(repository, lgr)
				// Add service middleware here
				// Logging middleware
				svc = middleware.LoggingMiddleware(lgr)(svc)
			}

			// Create Go kit endpoints for the Order Service
			// Then decorates with endpoint middlewares
			var endpoints usertansport.Endpoints
			{
				endpoints = usertansport.MakeEndpoints(svc)
			}
			// create http handler by endpoints
			var handler http.Handler
			{
				handler = httpusertransport.NewHandler(endpoints, lgr)
			}

			// initialize server config
			baseServer := &http.Server{
				Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
				Handler:      handler,
				ReadTimeout:  readTimeout,
				WriteTimeout: writeTimeout,
			}

			group, ctx := errgroup.WithContext(ctx)
			group.Go(func() error {
				if err := baseServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
					return fmt.Errorf("unable to start HTTP server: %w", err)
				}
				return nil
			})

			logger.Info(lgr, "HTTP server successfully started.")

			// Make sure to stop the http server when the context is canceled.
			group.Go(func() error {
				<-ctx.Done()
				logger.Info(lgr, "gracefully stopping HTTP server...")
				ctxShutdown, cancel := context.WithTimeout(ctx, shutdownTimeout)
				defer cancel()

				if err := baseServer.Shutdown(ctxShutdown); err != nil {
					return fmt.Errorf("HTTP server shutdown failed: %w", err)
				}

				logger.Info(lgr, "HTTP server gracefully stopped.")

				return nil
			})

			return group.Wait() //nolint:wrapcheck

			return nil
		},
	}
}
