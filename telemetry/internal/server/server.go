package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(ctx context.Context, log *slog.Logger, port string, register func(*grpc.Server)) error {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	if register != nil {
		register(grpcServer)
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Error("failed to listen:", slog.Any("error", err))
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)

	serverErrCh := make(chan error, 1)

	go func() {
		defer wg.Done()
		defer lis.Close()

		log.Info("gRPC server listening on " + port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("failed to serve gRPC server:", slog.Any("error", err))
			serverErrCh <- err
			return
		}
		serverErrCh <- nil
	}()

	<-ctx.Done()
	log.Info("shutting down gracefully, press Ctrl+C again to force")
	grpcServer.GracefulStop()
	wg.Wait()
	if err := <-serverErrCh; err != nil && !errors.Is(err, net.ErrClosed) {
		log.Error("gRPC serve returned error:", slog.Any("error", err))
	}

	log.Info("server stopped")
	return nil
}
