package main

import (
	"github.com/rs/zerolog/log"

	"github.com/888NiKiToS888/catalog-service/internal/app/config"
	rhealth "github.com/888NiKiToS888/catalog-service/internal/app/handler/http/health"
	rprocessor "github.com/888NiKiToS888/catalog-service/internal/app/processor/http"
)

func main() {
	config.Load()
	cfg := config.Root

	hHealth := rhealth.NewHandler()
	httpServer := rprocessor.NewHttp(hHealth, cfg.Processor.WebServer)
	if err := httpServer.Serve(); err != nil {
		log.Fatal().Err(err).Msg("HTTP server failed")
	}
}
