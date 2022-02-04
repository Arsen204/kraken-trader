package tg

import (
	"errors"
	"net/http"
	"net/url"
	"time"
)

//go:generate mockgen -source=tg.go -destination=mocks/tg.go

type TelegramConfig struct {
	Token  string
	ChatID string
}

type TelegramService interface {
	Send(text string) error
}

type telegramService struct {
	config TelegramConfig
}

func NewTelegramService(config TelegramConfig) TelegramService {
	return &telegramService{
		config: config,
	}
}

const clientRequestTimeout = 10 * time.Second

func (t *telegramService) Send(text string) error {
	endpoint := "https://api.telegram.org/bot" + t.config.Token + "/sendMessage"

	v := url.Values{}
	v.Add("chat_id", t.config.ChatID)
	v.Add("text", text)
	queryString := v.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint+"?"+queryString, nil)
	if err != nil {
		return err
	}

	client := http.Client{
		Timeout: clientRequestTimeout,
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		return errors.New(res.Status)
	}

	return nil
}
