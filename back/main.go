package main

import (
	"fmt"
	"gin-template/config"
	"gin-template/docs"
	"gin-template/logging"
	"gin-template/pkg/protocol/http"
	"os"
)

func setDoc(c config.Config) {
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Api Support"
	docs.SwaggerInfo.Description = "This is an API server."
	docs.SwaggerInfo.Version = c.Version
	if os.Getenv("BASE_PATH") != "" {
		docs.SwaggerInfo.Host = "localhost"
		docs.SwaggerInfo.BasePath = os.Getenv("BASE_PATH")
	} else {
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", "localhost", c.Server.Port)
		docs.SwaggerInfo.BasePath = "/api/v1"
	}
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  admin.admin@admin.fr

// @securityDefinitions.apikey  Bearer
// @in                          header
// @name                        Authorization
func main() {
	conf := config.NewConfig()
	setDoc(conf)
	logging.Error.Fatal(http.RunServer(conf))
}
