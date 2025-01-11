package interceptor

import (
	"context"
	"strings"

	"storya-gateway-backend/internal/client"
	"storya-gateway-backend/internal/pb/github.com/webbsalad/storya-passport-backend/passport"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var skipMethods = map[string]struct{}{
	"/otp.OtpService/SendOtp":    {},
	"/otp.OtpService/ConfirmOtp": {},

	"/passport.PassportService/RefreshToken": {},
	"/passport.PassportService/Register":     {},
	"/passport.PassportService/Login":        {},
}

func AuthInterceptor(grpcClients *client.GRPCClients) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if isAuthSkipped(method) {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		token, err := extractToken(ctx)
		if err != nil {
			return err
		}

		userID, deviceID, err := validateToken(ctx, grpcClients, token)
		if err != nil {
			return err
		}

		newCtx := attachMetadata(ctx, userID, deviceID)
		return invoker(newCtx, method, req, reply, cc, opts...)
	}
}

func isAuthSkipped(method string) bool {
	_, skip := skipMethods[method]
	return skip
}

func extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "metadata not found")
	}

	authHeaders := md.Get("Authorization")
	if len(authHeaders) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "missing authorization header")
	}

	token := strings.TrimPrefix(authHeaders[0], "Bearer ")
	if token == authHeaders[0] {
		return "", status.Errorf(codes.Unauthenticated, "invalid authorization token format")
	}

	return token, nil
}

func validateToken(ctx context.Context, grpcClients *client.GRPCClients, token string) (string, string, error) {
	resp, err := grpcClients.PassportClient.CheckToken(ctx, &passport.CheckTokenRequest{Token: token})
	if err != nil {
		return "", "", status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return resp.UserId, resp.DeviceId, nil
}

func attachMetadata(ctx context.Context, userID, deviceID string) context.Context {
	md := metadata.Pairs(
		"user_id", userID,
		"device_id", deviceID,
	)
	return metadata.NewOutgoingContext(ctx, md)
}
