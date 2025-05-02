package main

import (
	"github.com/shinplay/internal/api/router"
	"github.com/shinplay/internal/config"
)

func main() {
	appConfig := config.GetConfig()

	app := router.CreateNewFiberApp(appConfig)

	router.Routes(app)
	router.StartServer(app, appConfig)
}
