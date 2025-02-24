package mocks

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/mock"
)

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	args := m.Called(exchange, key, mandatory, immediate, msg)
	return args.Error(0)
}
