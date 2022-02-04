package rest

import (
	"course-project/internal/domain"
	"course-project/internal/service"

	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	serverCtx     context.Context
	serverStopCtx context.CancelFunc
	logger        logrus.FieldLogger
	services      *service.Service
}

func NewHandler(serverCtx context.Context, serverStopCtx context.CancelFunc, logger logrus.FieldLogger, services *service.Service) *Handler {
	return &Handler{
		serverCtx:     serverCtx,
		serverStopCtx: serverStopCtx,
		logger:        logger,
		services:      services,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/run", h.Run)
	r.Get("/stop", h.Stop)

	return r
}

func (h *Handler) Run(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("productID")
	if productID == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	size := r.URL.Query().Get("size")
	sz, err := strconv.ParseFloat(size, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	limitCf := r.URL.Query().Get("limitCf")
	lc, err := strconv.ParseFloat(limitCf, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	h.services.Algorithm.RunStochasticCross(h.serverCtx, domain.CandlePeriod(period), productID, sz, lc)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Stop(w http.ResponseWriter, r *http.Request) {
	h.serverStopCtx()
	w.WriteHeader(http.StatusOK)
}
