package amqprpc

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"

	"github.com/PanziApp/backend/internal/domain"
	"github.com/PanziApp/backend/internal/usecase"
	"github.com/PanziApp/backend/pkg/rabbitmq/rmq_rpc/server"
)

type translationRoutes struct {
	translationUseCase usecase.Translation
}

func newTranslationRoutes(routes map[string]server.CallHandler, t usecase.Translation) {
	r := &translationRoutes{t}
	{
		routes["getHistory"] = r.getHistory()
	}
}

type historyResponse struct {
	History []domain.Translation `json:"history"`
}

func (r *translationRoutes) getHistory() server.CallHandler {
	return func(d *amqp.Delivery) (interface{}, error) {
		translations, err := r.translationUseCase.History(context.Background())
		if err != nil {
			return nil, fmt.Errorf("amqp_rpc - translationRoutes - getHistory - r.translationUseCase.History: %w", err)
		}

		response := historyResponse{translations}

		return response, nil
	}
}
