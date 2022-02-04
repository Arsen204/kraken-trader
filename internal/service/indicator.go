package service

import (
	"course-project/internal/domain"

	"container/ring"
	"strconv"
)

//go:generate mockgen -source=indicator.go -destination=mocks/indicator.go

type IndicatorService interface {
	ClosePriceFilter(candles <-chan *domain.Candle) <-chan float64
	EMA(prices <-chan float64, n int) <-chan float64
	Stochastic(prices <-chan float64, n int) <-chan *StochasticStruct
}

type indicatorService struct{}

func NewIndicatorService() IndicatorService {
	return &indicatorService{}
}

func (i *indicatorService) ClosePriceFilter(candles <-chan *domain.Candle) <-chan float64 {
	ch := make(chan float64)

	go func() {
		defer close(ch)

		for candle := range candles {
			price, err := strconv.ParseFloat(candle.Close, 64)
			if price > 0 && err == nil {
				ch <- price
			}
		}
	}()

	return ch
}

func (i *indicatorService) EMA(prices <-chan float64, n int) <-chan float64 {
	output := make(chan float64)
	alpha := 2 / (float64(n) + 1)

	go func() {
		x1, x2 := 0.0, <-prices
		for price := range prices {
			x1 = x2
			x2 = alpha*price + (1-alpha)*x1
			output <- x2
		}
		close(output)
	}()

	return output
}

type StochasticStruct struct {
	K float64
	D float64
}

func (i *indicatorService) Stochastic(prices <-chan float64, n int) <-chan *StochasticStruct {
	output := make(chan *StochasticStruct)

	go func() {
		var min, max float64
		r := ring.New(n)

		sendEMA := make(chan float64)
		getEMA := i.EMA(sendEMA, 3)

		for k := 0; k < n; k++ {
			r.Value = <-prices
			r = r.Move(1)
		}

		sendEMA <- r.Prev().Value.(float64)

		for price := range prices {
			r.Value = price
			max = price
			min = max

			r.Do(func(i interface{}) {
				elem := i.(float64)
				if elem > max {
					max = elem
				}
				if elem < min {
					min = elem
				}
			})

			K := (price - min) / (max - min) * 100
			sendEMA <- K
			D := <-getEMA

			output <- &StochasticStruct{
				K: K,
				D: D,
			}
			r.Move(1)
		}
		close(sendEMA)
		close(output)
	}()

	return output
}
