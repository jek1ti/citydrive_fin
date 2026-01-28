package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jekiti/citydrive/api-gateway/internal/config"
	adminpb "github.com/jekiti/citydrive/gen/proto/admin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type AdminClient struct {
	client adminpb.AdminServiceClient
	conn   *grpc.ClientConn
}

func NewAdminClient(cfg *config.GRPCConfig) (*AdminClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		cfg.AdminAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to admin service at %s: %w", cfg.AdminAddr, err)
	}

	client := adminpb.NewAdminServiceClient(conn)

	log.Printf("Successfully connected to admin service at %s", cfg.AdminAddr)
	return &AdminClient{client: client, conn: conn}, nil
}



func (c *AdminClient) GetCar(ctx context.Context, traceID string, req *adminpb.GetCarRequest) (*adminpb.GetCarResponse, error) {
	if traceID == "" {
		return nil, fmt.Errorf("traceID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	md := metadata.Pairs("x-trace-id", traceID)
	ctx = metadata.NewOutgoingContext(ctx, md)

	
	response, err := c.client.GetCar(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}
	return response, nil
}

func (c *AdminClient) GetCarHistory(ctx context.Context, traceID string, req *adminpb.GetCarHistoryRequest) (*adminpb.GetCarHistoryResponse, error) {
	if traceID == "" {
		return nil, fmt.Errorf("traceID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	md := metadata.Pairs("x-trace-id", traceID)
	ctx = metadata.NewOutgoingContext(ctx, md)

	
	response, err := c.client.GetCarHistory(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to GetCarHistory: %w", err)
	}
	return response, nil
}

func (c *AdminClient) GetCarsHistory(ctx context.Context, traceID string, req *adminpb.GetCarsHistoryRequest) (*adminpb.GetCarsHistoryResponse, error) {
	if traceID == "" {
		return nil, fmt.Errorf("traceID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	md := metadata.Pairs("x-trace-id", traceID)
	ctx = metadata.NewOutgoingContext(ctx, md)

	
	response, err := c.client.GetCarsHistory(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to GetCarsHistory: %w", err)
	}
	return response, nil
}

func (c *AdminClient) GetCarsNow(ctx context.Context, traceID string, req *adminpb.GetCarsNowRequest) (*adminpb.GetCarsNowResponse, error) {
	if traceID == "" {
		return nil, fmt.Errorf("traceID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	md := metadata.Pairs("x-trace-id", traceID)
	ctx = metadata.NewOutgoingContext(ctx, md)

	
	response, err := c.client.GetCarsNow(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to GetCarsNow: %w", err)
	}
	return response, nil
}

func (c *AdminClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}