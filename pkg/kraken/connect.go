package kraken

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ErrLostConnection     = errors.New("lost connection with Kraken WebSocket")
	ErrKrakenConnect      = errors.New("cannot connect to Kraken WebSocket")
	ErrSubscriptionFailed = errors.New("subscription failed")
)

type (
	WSMessage struct {
		Event  string  `json:"event"`
		Candle *Candle `json:"candle"`
	}

	Candle struct {
		Open   string  `json:"open"`
		High   string  `json:"high"`
		Low    string  `json:"low"`
		Close  string  `json:"close"`
		Time   float64 `json:"time"`
		Volume float64 `json:"volume"`
	}
)

var candles chan []byte

func (k *krakenService) Candles(ctx context.Context, productID string, period string) (<-chan []byte, error) {
	err := k.connect(ctx, productID, period)
	if err != nil {
		return nil, err
	}
	return candles, nil
}

const pingPeriod = 70 * time.Second

func (k *krakenService) connect(ctx context.Context, productID string, period string) error {
	conn, err := k.openConnection()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrKrakenConnect, err)
	}

	err = subscribeToCandle(conn, productID, period)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSubscriptionFailed, err)
	}

	candles = make(chan []byte)

	go func(context.Context) {
		ticker := time.NewTicker(pingPeriod)

		defer func() {
			ticker.Stop()
			conn.Close()

			if r := recover(); r != nil {
				k.logger.Errorf("%v: %v", ErrLostConnection, r)

				err = k.connect(ctx, period, productID)
				if err != nil {
					k.logger.Error(err)
					close(candles)
				}
			} else {
				close(candles)
			}
		}()

		var recentTime float64

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				k.logger.Debug("ping connect")
				err := conn.WriteMessage(websocket.PingMessage, nil)
				if err != nil {
					k.logger.Error("ping message error")
					return
				}

			default:
				var message WSMessage

				err = conn.ReadJSON(&message)
				if err != nil {
					if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						panic(err)
					} else {
						k.logger.Error(err)
						return
					}
				}

				if message.Candle != nil && recentTime < message.Candle.Time {
					// k.logger.Debug(message.Candle)
					b, err := json.Marshal(message.Candle)
					if err != nil {
						k.logger.Error("cannot marshal candle")
					}

					candles <- b
					recentTime = message.Candle.Time
				}
			}
		}
	}(ctx)

	return nil
}

const (
	openConnectionTimeout = 10 * time.Second
	reconnectionTimeout   = 2 * time.Second
)

func (k *krakenService) openConnection() (*websocket.Conn, error) {
	u := url.URL{Scheme: "wss", Host: k.config.BaseURI, Path: k.config.WSEndpoint}
	k.logger.Infof("connecting to %s", u.String())

	ctx, ctxCancel := context.WithTimeout(context.Background(), openConnectionTimeout)
	defer ctxCancel()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				k.logger.Infof("reconnecting: %v", err)
				time.Sleep(reconnectionTimeout)
				continue
			}

			k.logger.Infof("connected to Kraken WebSocket")
			return conn, nil
		}
	}
}

func subscribeToCandle(conn *websocket.Conn, productID string, period string) error {
	subscriptionMsg := struct {
		Event      string   `json:"event"`
		Feed       string   `json:"feed"`
		ProductIDs []string `json:"product_ids"`
	}{
		Event:      "subscribe",
		Feed:       period,
		ProductIDs: []string{productID},
	}

	err := conn.WriteJSON(subscriptionMsg)
	if err != nil {
		return err
	}

	var msg WSMessage

	// Read first info message
	err = conn.ReadJSON(&msg)
	if err != nil {
		return err
	}

	// Read subscription status message
	err = conn.ReadJSON(&msg)
	if err != nil {
		return err
	}

	if msg.Event == "alert" {
		return errors.New("invalid productID or period")
	}

	return nil
}
