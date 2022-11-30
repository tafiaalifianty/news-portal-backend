package services

import (
	"errors"

	errn "final-project-backend/internal/errors"
	"final-project-backend/internal/helpers"
	"final-project-backend/internal/models"
	"final-project-backend/internal/repositories"

	"gorm.io/gorm"
)

type IAuthService interface {
	Register(user *models.User, referredCode *string) (*models.User, *string, *string, error)
	Login(email, password string) (*models.User, *string, *string, error)
	RefreshAccessToken(userID int64, token *string) (*string, error)
}

type authService struct {
	userRepository      repositories.IUserRepository
	userTokenRepository repositories.IUserTokenRepository
	hasher              helpers.Hasher
}

type AuthServiceConfig struct {
	userRepository      repositories.IUserRepository
	userTokenRepository repositories.IUserTokenRepository
	hasher              helpers.Hasher
}

func NewAuthService(c *AuthServiceConfig) IAuthService {
	return &authService{
		userRepository:      c.userRepository,
		userTokenRepository: c.userTokenRepository,
		hasher:              c.hasher,
	}
}

func (s *authService) Register(user *models.User, referredCode *string) (*models.User, *string, *string, error) {
	encryptedPassword, err := s.hasher.HashAndSalt(user.Password)
	if err != nil {
		return nil, nil, nil, err
	}

	user.Password = encryptedPassword
	user.Role = models.Member

	if *referredCode != "" {
		referredUser, err := s.userRepository.GetByReferralCode(*referredCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil, nil, errn.ErrInvalidReferralCode
			}
		}
		user.ReferredUserID = referredUser.ID
	}

	if user.Role == models.Member {
		user.ReferralCode = helpers.RandSeq(6)
	}

	createdUser, err := s.userRepository.Insert(user)
	if err != nil {
		return nil, nil, nil, err
	}

	dataToAccessToken := models.User{
		ID:    createdUser.ID,
		Email: createdUser.Email,
		Role:  createdUser.Role,
	}

	accessToken, err := helpers.GenerateAccessToken(dataToAccessToken)
	if err != nil {
		return nil, nil, nil, err
	}

	dataToRefreshToken := models.User{
		ID: createdUser.ID,
	}

	refreshToken, err := helpers.GenerateRefreshToken(dataToRefreshToken)
	if err != nil {
		return nil, nil, nil, err
	}

	_, err = s.userTokenRepository.Insert(&models.UserToken{
		UserID:       createdUser.ID,
		Email:        createdUser.Email,
		Role:         createdUser.Role,
		RefreshToken: *refreshToken,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return createdUser, accessToken, refreshToken, nil
}

func (s *authService) Login(email, password string) (*models.User, *string, *string, error) {
	user, err := s.userRepository.GetByEmail(email)
	if err != nil {
		return nil, nil, nil, err
	}

	if !s.hasher.ComparePasswords(user.Password, []byte(password)) {
		return nil, nil, nil, gorm.ErrRecordNotFound
	}

	dataToAccessToken := models.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}

	accessToken, err := helpers.GenerateAccessToken(dataToAccessToken)
	if err != nil {
		return nil, nil, nil, err
	}

	dataToRefreshToken := models.User{
		ID: user.ID,
	}

	refreshToken, err := helpers.GenerateRefreshToken(dataToRefreshToken)
	if err != nil {
		return nil, nil, nil, err
	}

	_, rowsAffected, err := s.userTokenRepository.Update(&models.UserToken{
		UserID:       user.ID,
		RefreshToken: *refreshToken,
	})
	if rowsAffected == 0 {
		_, err = s.userTokenRepository.Insert(&models.UserToken{
			UserID:       user.ID,
			Email:        user.Email,
			Role:         user.Role,
			RefreshToken: *refreshToken,
		})
	}

	if err != nil {
		return nil, nil, nil, err
	}

	return user, accessToken, refreshToken, nil
}

func (s *authService) RefreshAccessToken(userID int64, token *string) (*string, error) {
	userToken, err := s.userTokenRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if userToken.RefreshToken != *token {
		return nil, errn.ErrNotAuthorized
	}

	dataToAccessToken := models.User{
		ID:    userToken.UserID,
		Email: userToken.Email,
		Role:  userToken.Role,
	}

	accessToken, err := helpers.GenerateAccessToken(dataToAccessToken)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}
