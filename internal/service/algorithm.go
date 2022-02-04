package service

import (
	"course-project/internal/domain"

	"context"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=algorithm.go -destination=mocks/algorithm.go

type AlgorithmService interface {
	RunStochasticCross(ctx context.Context, period domain.CandlePeriod, symbol string, size float64, limitCf float64)
}

type algorithmService struct {
	logger    logrus.FieldLogger
	exchange  ExchangeService
	indicator IndicatorService
}

func NewAlgorithmService(logger logrus.FieldLogger, exchange ExchangeService, indicator IndicatorService) AlgorithmService {
	return &algorithmService{
		logger:    logger,
		exchange:  exchange,
		indicator: indicator,
	}
}

func (a *algorithmService) RunStochasticCross(ctx context.Context, period domain.CandlePeriod, symbol string, size float64, limitCf float64) {
	a.logger.Debug("RunStochasticCross")

	candles, err := a.exchange.Candles(ctx, symbol, period)
	if err != nil {
		a.logger.Error(err)
		return
	}
	prices := a.indicator.ClosePriceFilter(candles)

	go func() {
		currentPrice := <-prices
		highLimit := float64(int(currentPrice * (1 + limitCf)))
		lowLimit := float64(int(currentPrice * (1 - limitCf)))

		data := a.indicator.Stochastic(prices, 2)
		s := <-data

		var havePosition bool
		KIsHigherD := s.K > s.D

		for stochastic := range data {
			a.logger.Debug(stochastic)

			if stochastic.K > stochastic.D && !KIsHigherD {
				KIsHigherD = true

				if !havePosition {
					err := a.exchange.PlaceOrder(domain.OrderRequest{
						Symbol:     symbol,
						Side:       "buy",
						Size:       size,
						LimitPrice: highLimit,
					})

					if err != nil {
						a.logger.Error(err)
						return
					}

					havePosition = true
				}
			} else if stochastic.K < stochastic.D && KIsHigherD {
				KIsHigherD = false

				if havePosition {
					err := a.exchange.PlaceOrder(domain.OrderRequest{
						Symbol:     symbol,
						Side:       "sell",
						Size:       size,
						LimitPrice: lowLimit,
					})

					if err != nil {
						a.logger.Error(err)
						return
					}

					havePosition = false
				}
			}
		}
	}()
}
