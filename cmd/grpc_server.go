package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	"net"
	"os"
	"os/signal"
	"syscall"
	user "user-management/internal"
	"user-management/internal/config"
	"user-management/internal/config/db/postgress"
	"user-management/internal/dao/repo"
	userServ "user-management/internal/implementation"
	"user-management/internal/middleware"
	usertansport "user-management/internal/transport"
	grpchandler "user-management/internal/transport/grpc"
)

func newGRPCServerCmd(cfg config.Config, lgr log.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "grpc-server",
		Short: "Run app gRPC server",
		RunE: func(_ *cobra.Command, args []string) error {
			log.With(lgr, "newGRPCServerCommand")
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
			defer cancel()
			// DB creation
			var db *gorm.DB
			{
				var err error
				// Connect to the "ordersdb" database
				db, err = postgress.InitDB(cfg.DB, lgr)
				if err != nil {
					level.Error(lgr).Log("exit", err)
					os.Exit(-1)
				}
			}

			// Create Order Service
			var svc user.Service
			{
				repository, err := repo.NewUserSQLRepository(db, lgr)
				if err != nil {
					level.Error(lgr).Log("exit", err)
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
			var handler grpchandler.UserServer
			{
				handler = grpchandler.NewHandler(endpoints, lgr)
			}

			grpcServer := grpc.NewServer()
			reflection.Register(grpcServer)

			grpchandler.RegisterUserServer(grpcServer, handler)

			grpcListener, _ := net.Listen("tcp", ":9090")
			defer grpcListener.Close()

			group, ctx := errgroup.WithContext(ctx)

			group.Go(func() error {
				if err := grpcServer.Serve(grpcListener); err != nil {
					return fmt.Errorf("unable to start gRPC server: %w", err)
				}

				level.Info(lgr).Log("transport", "gRPC", "addr", "9090")
				return nil
			})

			group.Go(func() error {
				<-ctx.Done()
				level.Info(lgr).Log(ctx, "gracefully stopping gRPC server...")
				grpcServer.GracefulStop()

				return nil
			})

			return group.Wait()

		},
	}
}
