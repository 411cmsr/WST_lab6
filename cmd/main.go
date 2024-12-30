package main

import (
	"WST_lab6_server/config"

	"WST_lab6_server/internal/database/postgres"
	"WST_lab6_server/internal/httpserver/routes"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()
	postgres.Init()
	httpServer := gin.Default()

	routes.Init(httpServer)

	httpServer.StaticFile("/favicon.ico", "./favicon.ico")
	err := httpServer.Run(":8095")
	if err != nil {
		fmt.Println(err)
		return
	}
}
