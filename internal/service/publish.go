package service

import (
	"course-project/internal/domain"
	"course-project/internal/repository"
	tg "course-project/pkg/telegram"

	"context"
	"encoding/json"
	"io"

	"github.com/sirupsen/logrus"
)

type publishService struct {
	logger logrus.FieldLogger
	repo   repository.Repository
	tg     tg.TelegramService
}

func NewPublishService(logger logrus.FieldLogger, repo repository.Repository, tg tg.TelegramService) io.Writer {
	return &publishService{
		logger: logger,
		repo:   repo,
		tg:     tg,
	}
}

func (ps *publishService) Write(p []byte) (int, error) {
	var orderResponse domain.OrderResponse

	err := json.Unmarshal(p, &orderResponse)
	if err != nil {
		return 0, err
	}

	text, err := orderResponse.ToText()
	if err == nil {
		err := ps.tg.Send(text)
		if err != nil {
			ps.logger.Error("cannot send msg to telegram: ", err)
		}
	}

	err = ps.repo.CreateOrder(context.Background(), orderResponse.SendStatus.OrderEvents[0].OrderPriorExecution)
	if err != nil {
		return 0, err
	}

	return 0, nil
}
