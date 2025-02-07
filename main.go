package main

import (
	"log"

	_ "github.com/LeandroAlcantara-1997/investment-analyzer/docs"

	"github.com/LeandroAlcantara-1997/investment-analyzer/config/env"
	"github.com/LeandroAlcantara-1997/investment-analyzer/internal/app/container"
	"github.com/LeandroAlcantara-1997/investment-analyzer/internal/app/transport/http"
)

// @description     Investiment Analyzer investment analyzer is a project created to facilitate investment reports.
// @termsOfService  http://swagger.io/terms/

// @contact.url    https://www.linkedin.com/in/leandro-alcantara-pro

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
func main() {
	ctx, cont, err := container.New()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	http.New(env.Env.APIPort, env.Env.APIName,
		env.Env.APIVersion, env.Env.AllowOrigins,
		env.Env.Environment, cont).NewServer(ctx)
}
