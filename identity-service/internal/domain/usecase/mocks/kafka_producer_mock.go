package mocks

import (
	"github.com/stretchr/testify/mock"
)

type ProducerMock struct {
	mock.Mock
}

func (m *ProducerMock) SendUserCreated(userUUID string, email string) error {
	args := m.Called(userUUID, email)
	return args.Error(0)
}

func (m *ProducerMock) SendUserLogin(userUUID string, email string) error {
	args := m.Called(userUUID, email)
	return args.Error(0)
}
