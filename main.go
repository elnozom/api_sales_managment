package main

import (
	"hand_held/db"
	"hand_held/handler"
	"hand_held/router"
)

// @title Swagger Example API
// @version 1.0
// @description Conduit API
// @title Conduit API

// @host 127.0.0.1:8585
// @BasePath /api

// @schemes http https
// @produce	application/json
// @consumes application/json

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	r := router.New()

	v1 := r.Group("/api")
	db.InitDatabase()
	db := db.DBConn

	h, err := handler.NewHandler(db)
	if err != nil {
		panic(err)

	}
	h.Register(v1)
	r.Logger.Fatal(r.Start(":8586"))
}
