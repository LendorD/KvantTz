package mocks

import (
	"KvantTZ/internal/models"
	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *UserRepository) GetAll(offset, limit, minAge, maxAge int) ([]models.User, int64, error) {
	args := m.Called(offset, limit, minAge, maxAge)
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func (m *UserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepository) FindByID(id int) (*models.User, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepository) EmailExists(email string) bool {
	args := m.Called(email)
	return args.Bool(0)
}
