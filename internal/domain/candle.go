package domain

type CandlePeriod string

const (
	CandlePeriod1m  CandlePeriod = "candles_trade_1m"
	CandlePeriod5m  CandlePeriod = "candles_trade_5m"
	CandlePeriod15m CandlePeriod = "candles_trade_15m"
	CandlePeriod30m CandlePeriod = "candles_trade_30m"

	CandlePeriod1h  CandlePeriod = "candles_trade_1h"
	CandlePeriod4h  CandlePeriod = "candles_trade_4h"
	CandlePeriod12h CandlePeriod = "candles_trade_12h"

	CandlePeriod1d CandlePeriod = "candles_trade_1d"
	CandlePeriod1w CandlePeriod = "candles_trade_1w"
)

type Candle struct {
	Open   string  `json:"open"`
	High   string  `json:"high"`
	Low    string  `json:"low"`
	Close  string  `json:"close"`
	Time   float64 `json:"time"`
	Volume float64 `json:"volume"`
}
