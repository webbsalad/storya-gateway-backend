package docs

import (
	"log"
	"os"

	"github.com/swaggo/swag"
)

var swaggerInfo = &swag.Spec{
	Version:     "0.1.0",
	Title:       "Storya",
	Description: "GRPC-Gateway server",
}

func init() {
	data, err := os.ReadFile("internal/docs/gateway.swagger.json")
	if err != nil {
		log.Fatalf("failed to read swagger file: %v", err)
	}

	swaggerInfo.SwaggerTemplate = string(data)
	swag.Register("swagger", swaggerInfo)
}
