package main

import (
	"github.com/shinplay/internal/api/router"
)

func main() {
	// ctx := context.Background()
	app := router.CreateNewFiberApp()

	app.Routes()
	app.StartServer()
}
