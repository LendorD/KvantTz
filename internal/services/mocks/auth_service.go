package mocks

import (
	"github.com/stretchr/testify/mock"
)

type AuthService struct {
	mock.Mock
}

func NewAuthServiceMock() *AuthService {
	return &AuthService{}
}

func (m *AuthService) Login(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}
