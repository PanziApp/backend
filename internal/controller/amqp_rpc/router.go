package amqprpc

import (
	"github.com/PanziApp/backend/internal/usecase"
	"github.com/PanziApp/backend/pkg/rabbitmq/rmq_rpc/server"
)

// NewRouter -.
func NewRouter(t usecase.Translation) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)
	{
		newTranslationRoutes(routes, t)
	}

	return routes
}
