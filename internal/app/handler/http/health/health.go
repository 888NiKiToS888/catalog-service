package rhealth

import (
	"net/http"

	"github.com/rs/zerolog/log"

	rhandler "github.com/888NiKiToS888/catalog-service/internal/app/handler/http"
)

type handler struct{}

func NewHandler() rhandler.Health {
	return &handler{}
}

func (h *handler) LastCheck(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Error().Err(err).Msg("failed to write health response")
	}
}
