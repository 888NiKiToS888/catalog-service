package main

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/888NiKiToS888/catalog-service/internal/app/config"
	hcategory "github.com/888NiKiToS888/catalog-service/internal/app/handler/http/category"
	rhealth "github.com/888NiKiToS888/catalog-service/internal/app/handler/http/health"
	hproduct "github.com/888NiKiToS888/catalog-service/internal/app/handler/http/product"
	rprocessor "github.com/888NiKiToS888/catalog-service/internal/app/processor/http"
	pcategory "github.com/888NiKiToS888/catalog-service/internal/app/repository/category"
	rcpostgres "github.com/888NiKiToS888/catalog-service/internal/app/repository/conn/postgres"
	pproduct "github.com/888NiKiToS888/catalog-service/internal/app/repository/product"
	scategory "github.com/888NiKiToS888/catalog-service/internal/app/service/category"
	sproduct "github.com/888NiKiToS888/catalog-service/internal/app/service/product"
)

func main() {
	config.Load()
	cfg := config.Root

	ctx := context.Background()

	// Подключение к PostgreSQL
	pgClient, err := rcpostgres.NewConn(ctx, cfg.Repository.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}

	// Миграции
	oldVer, newVer, err := pgClient.Migrate(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}
	if oldVer != newVer {
		log.Info().Int64("old_version", oldVer).Int64("new_version", newVer).Msg("Database migrated")
	} else {
		log.Info().Int64("version", newVer).Msg("Database is up to date")
	}

	// Репозитории
	categoryRepo := pcategory.NewRepoFromPostgres(pgClient)
	productRepo := pproduct.NewRepoFromPostgres(pgClient)

	// Сервисы
	categorySvc := scategory.NewService(categoryRepo, productRepo)
	productSvc := sproduct.NewService(productRepo, categoryRepo)

	// Хендлеры
	healthHandler := rhealth.NewHandler()
	categoryHandler := hcategory.NewHandler(categorySvc)
	productHandler := hproduct.NewHandler(productSvc)

	// HTTP-сервер
	httpServer := rprocessor.NewHttp(
		cfg.Processor.WebServer,
		healthHandler,
		categoryHandler,
		productHandler,
	)
	if err := httpServer.Serve(); err != nil {
		log.Fatal().Err(err).Msg("HTTP server failed")
	}
}
