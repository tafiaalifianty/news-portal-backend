package services

import (
	"final-project-backend/internal/models"
	"final-project-backend/internal/repositories"
)

type IHistoryService interface {
	UpdateOrInsert(history *models.History) (*models.History, error)
	GetByUserAndPostID(userID int64, postID int64) (*models.History, error)
	GetByUserID(userID int64) ([]*models.History, error)
}

type historyService struct {
	historyRepository repositories.IHistoryRepository
}

type HistoryServiceConfig struct {
	historyRepository repositories.IHistoryRepository
}

func NewHistoryService(c *HistoryServiceConfig) IHistoryService {
	return &historyService{
		historyRepository: c.historyRepository,
	}
}

func (s *historyService) UpdateOrInsert(history *models.History) (*models.History, error) {
	newHistory, rowsAffected, err := s.historyRepository.Update(history)

	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		newHistory, insertErr := s.historyRepository.Insert(history)

		if insertErr != nil {
			return nil, insertErr
		}

		return newHistory, nil
	}

	return newHistory, nil
}

func (s *historyService) GetByUserID(userID int64) ([]*models.History, error) {
	histories, err := s.historyRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return histories, nil
}

func (s *historyService) GetByUserAndPostID(userID int64, postID int64) (*models.History, error) {
	history, err := s.historyRepository.GetByUserAndPostID(userID, postID)
	if err != nil {
		return nil, err
	}

	return history, nil
}
