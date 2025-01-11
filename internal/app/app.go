package app

import (
	"storya-gateway-backend/internal/client"
	"storya-gateway-backend/internal/config"

	"go.uber.org/fx"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(
			config.NewConfig,
			client.NewGRPCClients,
		),
		routerOption(),
	)
}
