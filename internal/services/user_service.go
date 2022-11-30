package services

import (
	"final-project-backend/internal/models"
	"final-project-backend/internal/repositories"

	"gorm.io/gorm"
)

type IUserService interface {
	GetByID(id int64) (*models.User, error)
	GetUsersReferredTotalSpending(id int64) ([]*models.UserSpendingTotalAggregates, error)
	Update(user *models.User) (*models.User, error)
}

type userService struct {
	userRepository repositories.IUserRepository
}

type UserServiceConfig struct {
	userRepository repositories.IUserRepository
}

func NewUserService(c *UserServiceConfig) IUserService {
	return &userService{
		userRepository: c.userRepository,
	}
}

func (s *userService) GetByID(id int64) (*models.User, error) {
	result, err := s.userRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *userService) GetUsersReferredTotalSpending(id int64) ([]*models.UserSpendingTotalAggregates, error) {
	users, err := s.userRepository.GetUsersReferredTotalSpending(id)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *userService) Update(user *models.User) (*models.User, error) {
	result, rowsAffected, err := s.userRepository.Update(user)

	if rowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}
