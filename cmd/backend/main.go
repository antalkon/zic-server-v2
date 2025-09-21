package main

import (
	"backend/internal/app"
	"log"

	_ "github.com/swaggo/echo-swagger"
)

// @title           Zentas Informatics Class srv API
// @version         2.0
// @description     API документация для серверного клиента Zentas Informatics Class 2 версии (ревизия)

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// ===== Cookie-based auth =====

// @securityDefinitions.apikey AccessToken
// @type apiKey
// @in cookie
// @name access_token

// (опционально, если хочешь отдельно авторизовываться и по refresh)
// @securityDefinitions.apikey RefreshToken
// @type apiKey
// @in cookie
// @name refresh_token

// (если хочешь сохранить поддержку заголовка тоже — оставь)
// @securityDefinitions.apikey BearerAuth
// @type apiKey
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("Application initialization failed: %v", err)
	}

	if err := application.RunServer(); err != nil {
		log.Fatalf("❌ Failed to run server: %v", err)
	}
}
