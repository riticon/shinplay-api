package main

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/shinplay/ent"
	"github.com/shinplay/internal/api/router"
	"github.com/shinplay/internal/config"
)

func main() {
	appConfig := config.GetConfig()
	client, err := ent.Open("postgres", "host=127.0.0.1 port=5432 user=shinplay dbname=shinplay_development password=shinplay sslmode=disable")

	if err != nil {
		panic(err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		appConfig.Logger.Error("failed creating schema resources:")
	}
	ctx := context.Background()
	// u, err := client.User.Create().SetEmail("pavankumarv124@gmail.com").SetUsername("pavan").Save(ctx)
	// // if err != nil {
	// // 	print("Error creating user")
	// // }
	// print(u.ID)
	users, err := client.User.Query().All(ctx)

	for i := range users {
		user := users[i]
		fmt.Printf("user %d \n", user.ID)
	}
	app := router.CreateNewFiberApp(appConfig)

	app.Routes()
	app.StartServer()
}
