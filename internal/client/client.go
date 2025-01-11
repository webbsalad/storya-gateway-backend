package client

import (
	"fmt"
	"log"
	"storya-gateway-backend/internal/config"
	"storya-gateway-backend/internal/pb/github.com/webbsalad/storya-content-backend/content"
	"storya-gateway-backend/internal/pb/github.com/webbsalad/storya-otp-backend/otp"
	"storya-gateway-backend/internal/pb/github.com/webbsalad/storya-passport-backend/passport"
	"storya-gateway-backend/internal/pb/github.com/webbsalad/storya-recs-backend/recs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClients struct {
	OtpClient         otp.OtpServiceClient
	PassportClient    passport.PassportServiceClient
	ContentClient     content.ContentServiceClient
	UserContentClient content.UserContentServiceClient
	RecsClient        recs.RecsServiceClient
}

func NewGRPCClients(cfg config.Config) *GRPCClients {
	otpConn, err := newGRPCClientConn(cfg.OtpAddr)
	if err != nil {
		log.Fatalf("connecting otp service: %v", err)
	}

	passportConn, err := newGRPCClientConn(cfg.PassportAddr)
	if err != nil {
		log.Fatalf("connecting passport service: %v", err)
	}

	contentConn, err := newGRPCClientConn(cfg.ContentAddr)
	if err != nil {
		log.Fatalf("connecting content service: %v", err)
	}

	recsConn, err := newGRPCClientConn(cfg.RecsAddr)
	if err != nil {
		log.Fatalf("connect recs service: %v", err)
	}

	return &GRPCClients{
		OtpClient:      otp.NewOtpServiceClient(otpConn),
		PassportClient: passport.NewPassportServiceClient(passportConn),
		ContentClient:  content.NewContentServiceClient(contentConn),
		RecsClient:     recs.NewRecsServiceClient(recsConn),
	}
}

func newGRPCClientConn(addr string) (*grpc.ClientConn, error) {
	if addr == "" {
		return nil, fmt.Errorf("address is required")
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to service at %s: %w", addr, err)
	}

	return conn, nil
}
