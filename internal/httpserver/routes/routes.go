package routes

import (
	"WST_lab6_server/internal/database/postgres"
	"WST_lab6_server/internal/handlers"
	"WST_lab6_server/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Init(httpserver *gin.Engine) {
	//middleware для обработки ошибок
	httpserver.Use(middleware.ErrorHandler())
	//Восстановление после паники
	httpserver.Use(gin.Recovery())
	//Логгирование 
	httpserver.Use(gin.Logger())
	//Подключение к БД
	db := postgres.Init()
	storage := &postgres.Storage{DB: db}
	route := &handlers.StorageHandler{Storage: storage}
	//routes по запросам
	apiv1 := httpserver.Group("/api/v1")
	apiv1.GET("/persons", route.SearchPersonHandler)
	apiv1.POST("/persons", middleware.BasicAuthMiddleware(), route.AddPersonHandler)
	apiv1.GET("/persons/list", route.GetAllPersonsHandler)
	apiv1.GET("/person/:id", route.GetPersonHandler)
	apiv1.PUT("/person/:id",middleware.BasicAuthMiddleware(), route.UpdatePersonHandler)
	apiv1.DELETE("/person/:id", middleware.BasicAuthMiddleware(), route.DeletePersonHandler)

}
