package services

import (
	"final-project-backend/internal/constants"
	errn "final-project-backend/internal/errors"
	"final-project-backend/internal/models"
	"final-project-backend/internal/repositories"

	"gorm.io/gorm"
)

type IGiftService interface {
	GetAll() ([]*models.Gift, error)
	GetAllUserGifts() ([]*models.UserGift, error)
	GetUserGiftsByID(userID int64) ([]*models.UserGift, error)
	UpdateUserGiftStatus(*models.UserGift) (*models.UserGift, error)
	UpdateGiftStock(giftID int64, newStock int) (*models.Gift, error)
}

type giftService struct {
	giftRepository repositories.IGiftRepository
}

type GiftServiceConfig struct {
	giftRepository repositories.IGiftRepository
}

func NewGiftService(c *GiftServiceConfig) IGiftService {
	return &giftService{giftRepository: c.giftRepository}
}

func (s *giftService) GetAll() ([]*models.Gift, error) {
	Gifts, err := s.giftRepository.GetAll()

	if err != nil {
		return nil, err
	}

	return Gifts, nil
}

func (s *giftService) GetAllUserGifts() ([]*models.UserGift, error) {
	userGifts, err := s.giftRepository.GetAllUserGifts()

	if err != nil {
		return nil, err
	}

	return userGifts, nil
}

func (s *giftService) GetUserGiftsByID(userID int64) ([]*models.UserGift, error) {
	userGifts, err := s.giftRepository.GetUserGiftsByUserID(userID)

	if err != nil {
		return nil, err
	}

	return userGifts, nil
}

func (s *giftService) UpdateUserGiftStatus(userGift *models.UserGift) (*models.UserGift, error) {
	currentUserGift, err := s.giftRepository.GetUserGift(&models.UserGift{
		UserID: userGift.UserID,
		GiftID: userGift.GiftID,
		Month:  userGift.Month,
		Year:   userGift.Year,
	})
	if err != nil {
		return nil, err
	}

	if currentUserGift.Status != constants.PROCESSED {
		return nil, errn.ErrInvalidStatusUpdate
	}
	updatedUserGift, rowsAffected, err := s.giftRepository.UpdateUserGift(userGift)

	if rowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if err != nil {
		return nil, err
	}

	return updatedUserGift, nil
}

func (s *giftService) UpdateGiftStock(giftID int64, newStock int) (*models.Gift, error) {
	updatedGift, rowsAffected, err := s.giftRepository.UpdateStock(&models.Gift{
		Model: gorm.Model{
			ID: uint(giftID),
		},
		Stock: newStock,
	})

	if rowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if err != nil {
		return nil, err
	}

	return updatedGift, nil
}
