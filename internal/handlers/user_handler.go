package handlers

import (
	"errors"
	"net/http"

	"final-project-backend/internal/dtos"
	errn "final-project-backend/internal/errors"
	"final-project-backend/internal/helpers"
	"final-project-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) GetUserProfile(c *gin.Context) {
	userContext, ok := c.Get("user")
	if !ok {
		helpers.SendErrorResponse(
			c,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	var response dtos.GetUserProfileResponseDTO

	user, err := h.services.User.GetByID(userContext.(dtos.JwtData).ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.SendErrorResponse(c, http.StatusNotFound, errn.ErrUserNotFound.Error())
			return
		}

		helpers.SendErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	response = dtos.GetUserProfileResponseDTO{
		ID:           user.ID,
		Email:        user.Email,
		Role:         user.Role,
		Fullname:     user.Fullname,
		Address:      user.Address,
		ReferralCode: user.ReferralCode,
	}

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) GetUserHistories(c *gin.Context) {
	var response dtos.GetUserHistoriesDTO

	userContext, ok := c.Get("user")
	if !ok {
		helpers.SendErrorResponse(
			c,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	histories, err := h.services.History.GetByUserID(userContext.(dtos.JwtData).ID)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response = dtos.FormatHistories(histories)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) GetUserSubscriptions(c *gin.Context) {
	var response dtos.GetUserSubscriptionsDTO

	userContext, ok := c.Get("user")
	if !ok {
		helpers.SendErrorResponse(
			c,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	userSubscriptions, err := h.services.UserSubscription.GetAllUserSubscriptions(userContext.(dtos.JwtData).ID)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response = dtos.FormatUserSubscriptions(userSubscriptions)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) GetUserInvoices(c *gin.Context) {
	var response dtos.GetAllInvoicesResponse

	userContext, ok := c.Get("user")
	if !ok {
		helpers.SendErrorResponse(
			c,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	invoices, err := h.services.Invoice.GetUserInvoices(userContext.(dtos.JwtData).ID)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	response = dtos.FormatInvoices(invoices)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) GetUserGifts(c *gin.Context) {
	var response dtos.GetUserGiftsResponse

	userContext, ok := c.Get("user")
	if !ok {
		helpers.SendErrorResponse(
			c,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	userGifts, err := h.services.Gift.GetUserGiftsByID(userContext.(dtos.JwtData).ID)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	response = dtos.FormatUserGifts(userGifts)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) GetUserReferrals(c *gin.Context) {
	var response dtos.GetUsersReferralsResponse

	userContext, ok := c.Get("user")
	if !ok {
		helpers.SendErrorResponse(
			c,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	users, err := h.services.User.GetUsersReferredTotalSpending(userContext.(dtos.JwtData).ID)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	response = users

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) GetUserVouchers(c *gin.Context) {
	var response dtos.GetUserVouchersResponse

	userContext, ok := c.Get("user")
	if !ok {
		helpers.SendErrorResponse(
			c,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	userVouchers, err := h.services.Voucher.GetUserVouchersByID(userContext.(dtos.JwtData).ID)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	response = dtos.FormatUserVouchers(userVouchers)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	var request dtos.UpdateUserRequest
	var response dtos.UpdateUserResponse

	userContext, ok := c.Get("user")
	if !ok {
		helpers.SendErrorResponse(
			c,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	err := c.ShouldBindJSON(&request)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	user := &models.User{
		ID:       userContext.(dtos.JwtData).ID,
		Fullname: request.Fullname,
		Address:  request.Address,
	}

	updatedUser, err := h.services.User.Update(user)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response = *dtos.FormatUpdattedUser(updatedUser)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}
