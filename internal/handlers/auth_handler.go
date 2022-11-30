package handlers

import (
	"errors"
	"net/http"

	"final-project-backend/internal/dtos"
	errn "final-project-backend/internal/errors"
	"final-project-backend/internal/helpers"
	"final-project-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func (h *Handler) Register(c *gin.Context) {
	var request dtos.RegisterRequestDTO
	var response dtos.RegisterResponseDTO

	err := c.ShouldBindJSON(&request)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	user := &models.User{
		Fullname: request.Fullname,
		Password: request.Password,
		Email:    request.Email,
		Address:  request.Address,
	}

	createdUser, accessToken, refreshToken, err := h.services.Auth.Register(user, request.ReferralCode)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == errn.UniqueViolation {
			helpers.SendErrorResponse(c, http.StatusBadRequest, errn.ErrEmailAlreadyExist.Error())

			return
		}

		if errors.Is(err, errn.ErrInvalidReferralCode) {
			helpers.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		helpers.SendErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	response = dtos.RegisterResponseDTO{
		ID:           createdUser.ID,
		Fullname:     createdUser.Fullname,
		Email:        createdUser.Email,
		Address:      createdUser.Address,
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}

	helpers.SendSuccessResponse(c, http.StatusCreated, http.StatusText(http.StatusCreated), response)
}

func (h *Handler) Login(c *gin.Context) {
	var request dtos.LoginRequestDTO
	var response dtos.LoginResponseDTO

	err := c.ShouldBindJSON(&request)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	user, accessToken, refreshToken, err := h.services.Auth.Login(request.Email, request.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.SendErrorResponse(c, http.StatusUnauthorized, errn.ErrWrongEmailOrPassword.Error())
			return
		}

		helpers.SendErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	response = dtos.LoginResponseDTO{
		ID:           user.ID,
		Email:        user.Email,
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var response dtos.RefreshTokenResponse

	refreshTokenContext, ok := c.Get("refresh_token")
	if !ok {
		helpers.SendErrorResponse(
			c,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	refreshToken := refreshTokenContext.(dtos.RefreshTokenContext)

	accessToken, err := h.services.Auth.RefreshAccessToken(refreshToken.ID, &refreshToken.RefreshToken)
	if err != nil {
		helpers.SendErrorResponse(
			c,
			http.StatusUnauthorized,
			http.StatusText(http.StatusUnauthorized),
		)
		return
	}

	response = dtos.RefreshTokenResponse{
		AccessToken: *accessToken,
	}

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}
