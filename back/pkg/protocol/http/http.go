package http

import (
	"fmt"
	"gin-template/config"
	"gin-template/database"
	"gin-template/pkg/middleware"
	v1 "gin-template/pkg/service/v1"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RunServer(conf config.Config) error {
	db := database.NewDatabase(conf.Db)

	r := gin.New()
	dbM := middleware.DatabaseMiddleware{DB: db}
	r.Use(gin.LoggerWithFormatter(middleware.CustomLogger), gin.Recovery(), middleware.CORSMiddleware(), dbM.SetDBMiddleware())

	if conf.Env == config.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup the routes.
	rg := r.Group("/api/v1")
	// Setup the routes for the auth service.
	v1.SetAuthService(rg.Group("/auth"), conf.Jwt)
	// Setup the routes for the user service.
	v1.SetUserRoutes(rg.Group("/users"), conf.Jwt)

	if conf.Env == config.Development {
		r.GET("/doc/*any", ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.DefaultModelsExpandDepth(-1),
		))
	}
	return r.Run(fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port))
}
