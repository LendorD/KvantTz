package mocks

import (
	"KvantTZ/internal/models"
	"github.com/stretchr/testify/mock"
)

type UserService struct {
	mock.Mock
}

// CreateUser
func (m *UserService) CreateUser(req *models.CreateUserRequest) (*models.UserResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

// GetAllUsers
func (m *UserService) GetAllUsers(page, limit, minAge, maxAge int) ([]models.UserResponse, int64, error) {
	args := m.Called(page, limit, minAge, maxAge)
	return args.Get(0).([]models.UserResponse), args.Get(1).(int64), args.Error(2)
}

// GetUserByID
func (m *UserService) GetUserByID(id int) (*models.UserResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

// UpdateUser
func (m *UserService) UpdateUser(id int, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

// DeleteUser
func (m *UserService) DeleteUser(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
