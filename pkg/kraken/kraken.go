package kraken

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=kraken.go -destination=mocks/kraken.go

type (
	Service interface {
		CandleGetter
		OrderPlacer
	}

	CandleGetter interface {
		Candles(ctx context.Context, productID string, period string) (<-chan []byte, error)
	}

	OrderPlacer interface {
		PlaceOrder(order OrderRequest) error
	}
)

type Config struct {
	BaseURI      string
	BaseDemoURI  string
	WSEndpoint   string
	RestEndpoint string
	APIKey       string
	APISecret    string
}

type krakenService struct {
	logger    logrus.FieldLogger
	config    *Config
	publisher io.Writer
}

func NewKrakenService(logger logrus.FieldLogger, config Config, publisher io.Writer) Service {
	return &krakenService{
		logger:    logger,
		config:    &config,
		publisher: publisher,
	}
}
