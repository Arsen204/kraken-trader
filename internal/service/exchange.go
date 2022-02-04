package service

import (
	"course-project/internal/domain"
	"course-project/pkg/kraken"

	"context"
	"encoding/json"
	"io"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=exchange.go -destination=mocks/exchange.go

type ExchangeService interface {
	Candles(ctx context.Context, productID string, period domain.CandlePeriod) (<-chan *domain.Candle, error)
	PlaceOrder(order domain.OrderRequest) error
}

type exchangeService struct {
	kraken kraken.Service
}

func NewExchangeService(logger logrus.FieldLogger, config kraken.Config, publisher io.Writer) ExchangeService {
	return &exchangeService{
		kraken: kraken.NewKrakenService(logger, config, publisher),
	}
}

func (w *exchangeService) Candles(ctx context.Context, productID string, period domain.CandlePeriod) (<-chan *domain.Candle, error) {
	candles, err := w.kraken.Candles(ctx, productID, string(period))
	if err != nil {
		return nil, err
	}

	ch := make(chan *domain.Candle)
	go func() {
		defer close(ch)

		for candle := range candles {
			var cdl domain.Candle
			err := json.Unmarshal(candle, &cdl)
			if err == nil {
				ch <- &cdl
			}
		}
	}()

	return ch, nil
}

func (w *exchangeService) PlaceOrder(order domain.OrderRequest) error {
	return w.kraken.PlaceOrder(kraken.OrderRequest{
		Symbol:     order.Symbol,
		Side:       order.Side,
		Size:       order.Size,
		LimitPrice: order.LimitPrice,
	})
}
