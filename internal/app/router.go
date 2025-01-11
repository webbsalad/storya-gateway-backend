package app

import (
	"context"
	"log"
	"net/http"

	"storya-gateway-backend/internal/aggregator"
	"storya-gateway-backend/internal/app/interceptor"
	"storya-gateway-backend/internal/client"
	"storya-gateway-backend/internal/config"
	"storya-gateway-backend/internal/pb/github.com/webbsalad/storya-content-backend/content"
	"storya-gateway-backend/internal/pb/github.com/webbsalad/storya-otp-backend/otp"
	"storya-gateway-backend/internal/pb/github.com/webbsalad/storya-passport-backend/passport"
	"storya-gateway-backend/internal/pb/github.com/webbsalad/storya-recs-backend/recs"

	"go.uber.org/fx"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func routerOption() fx.Option {
	return fx.Options(
		fx.Invoke(
			func(lc fx.Lifecycle, cfg config.Config, grpcClients *client.GRPCClients) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go newRouter(cfg, grpcClients)
						return nil
					},
				})
			},
		),
	)
}

func permissions(cfg config.Config, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		origin := r.Header.Get("Origin")
		if isAllowedOrigin(cfg.AllowedOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType, Grpc-Metadata-Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func newRouter(cfg config.Config, grpcClients *client.GRPCClients) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gatewayMux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			interceptor.AuthInterceptor(grpcClients),
		),
	}

	registerServices(ctx, gatewayMux, cfg, opts)

	httpMux := http.NewServeMux()

	httpMux.Handle("/", gatewayMux)
	httpMux.HandleFunc("/docs/", httpSwagger.WrapHandler)

	httpMux.HandleFunc("/mixed", aggregator.MixedHandler(cfg, grpcClients))
	httpMux.HandleFunc("/server-mixed", aggregator.MixedClientsHandler(cfg, grpcClients))

	log.Printf("Server listening at :50060")
	if err := http.ListenAndServe(":50060", permissions(cfg, httpMux)); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func registerServices(ctx context.Context, gatewayMux *runtime.ServeMux, cfg config.Config, opts []grpc.DialOption) {
	if err := passport.RegisterPassportServiceHandlerFromEndpoint(ctx, gatewayMux, cfg.PassportAddr, opts); err != nil {
		log.Fatalf("Failed to register Passport service: %v", err)
	}

	if err := otp.RegisterOtpServiceHandlerFromEndpoint(ctx, gatewayMux, cfg.OtpAddr, opts); err != nil {
		log.Fatalf("Failed to register OTP service: %v", err)
	}

	if err := content.RegisterContentServiceHandlerFromEndpoint(ctx, gatewayMux, cfg.ContentAddr, opts); err != nil {
		log.Fatalf("Failed to register OTP service: %v", err)
	}

	if err := content.RegisterUserContentServiceHandlerFromEndpoint(ctx, gatewayMux, cfg.ContentAddr, opts); err != nil {
		log.Fatalf("Failed to register OTP service: %v", err)
	}

	if err := recs.RegisterRecsServiceHandlerFromEndpoint(ctx, gatewayMux, cfg.RecsAddr, opts); err != nil {
		log.Fatalf("Failed to register OTP service: %v", err)
	}
}

func isAllowedOrigin(allowedOrigins []string, origin string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}
