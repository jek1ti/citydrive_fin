package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jekiti/citydrive/api-gateway/internal/config"
	authpb "github.com/jekiti/citydrive/gen/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type AuthClient struct {
	client authpb.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthClient(cfg *config.GRPCConfig) (*AuthClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		cfg.AuthAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service at %s: %w", cfg.AuthAddr, err)
	}

	client := authpb.NewAuthServiceClient(conn)

	log.Printf("Successfully connected to auth service at %s", cfg.AuthAddr)
	return &AuthClient{client: client, conn: conn}, nil
}

func (c *AuthClient) Login(ctx context.Context, traceID string, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	if traceID == "" {
		return nil, fmt.Errorf("traceID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	md := metadata.Pairs("x-trace-id", traceID)
	ctx = metadata.NewOutgoingContext(ctx, md)

	if req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("email and password cannot be empty")
	}
	response, err := c.client.Login(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}
	return response, nil
}

func (c *AuthClient) Register(ctx context.Context, traceID string, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	if traceID == "" {
		return nil, fmt.Errorf("traceID cannot be empty")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	md := metadata.Pairs("x-trace-id", traceID)
	ctx = metadata.NewOutgoingContext(ctx, md)

	if req.Email == "" || req.Name == "" || req.Surname == "" || req.Department == "" {
		return nil, fmt.Errorf("email, name, surname, and department cannot be empty")
	}
	response, err := c.client.Register(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return response, nil
}

func (c *AuthClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
