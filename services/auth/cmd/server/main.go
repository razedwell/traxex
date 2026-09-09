package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authv1 "github.com/razedwell/traxex/proto/gen/go/auth/v1"
	"github.com/razedwell/traxex/services/auth/internal/handler"
	"github.com/razedwell/traxex/services/auth/internal/handler/authhttp"
	"github.com/razedwell/traxex/services/auth/internal/repository"
	"github.com/razedwell/traxex/services/auth/internal/service"
	"github.com/razedwell/traxex/shared/config"
	"github.com/razedwell/traxex/shared/obs"
	"github.com/razedwell/traxex/shared/postgres"
	"github.com/razedwell/traxex/shared/redis"
	"github.com/razedwell/traxex/shared/traxerr"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger, err := obs.InitLogger(cfg.Logging)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	tracerShutdown, err := obs.InitTracer(ctx, "auth-service")
	if err != nil {
		logger.Fatal("failed initializing tracer", zap.Error(err))
	}

	defer func() {
		sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tracerShutdown(sctx); err != nil {
			logger.Error("tracer shutdown failed", zap.Error(err))
		}
	}()

	pool, err := postgres.NewPool(ctx, cfg.Database)
	if err != nil {
		logger.Fatal("failed to initialize pg pool", zap.Error(err))
	}
	defer pool.Close()

	rdb, err := redis.NewClient(ctx, cfg.Redis)
	if err != nil {
		logger.Fatal("failed to initialize redis client", zap.Error(err))
	}
	defer rdb.Close()

	//object graph bottom-up
	userRepo := repository.NewPostgresUserRepo(pool)
	rdsesh := repository.NewRedisSessionStore(rdb)
	authsvc := service.NewAuthService(userRepo, rdsesh, []byte(cfg.Auth.JWTSecret))
	grpcHandler := handler.NewGRPCHandler(authsvc)
	httpHandler := authhttp.NewHandler(authsvc)

	//grpc server
	gs := grpc.NewServer()
	authv1.RegisterAuthServiceServer(gs, grpcHandler)
	reflection.Register(gs)

	//http server
	strict := authhttp.NewStrictHandlerWithOptions(httpHandler, nil, authhttp.StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			var te *traxerr.Traxerr
			if errors.As(err, &te) {
				http.Error(w, te.Message, traxerr.HTTPStatus(te.Code))
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		},
	})
	hs := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: authhttp.Handler(strict),
	}

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Auth.GRPCPort))
		if err != nil {
			return fmt.Errorf("error listening on grpc port: %w", err)
		}
		logger.Info("grpc server listening", zap.Int("port", cfg.Auth.GRPCPort))
		return gs.Serve(lis)
	})

	g.Go(func() error {
		logger.Info("http server listening", zap.Int("port", cfg.Server.Port))
		if err := hs.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-gctx.Done()
		logger.Info("shutting down")
		gs.GracefulStop()
		sctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return hs.Shutdown(sctx)
	})

	if err := g.Wait(); err != nil {
		logger.Fatal("server exited", zap.Error(err))
	}
	logger.Info("shutdown complete")
}
