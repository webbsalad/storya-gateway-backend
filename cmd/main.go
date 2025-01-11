package main

import (
	"storya-gateway-backend/internal/app"
	_ "storya-gateway-backend/internal/docs"
)

func main() {
	app.NewApp().Run()
}
