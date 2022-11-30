package handlers

import (
	"net/http"

	"final-project-backend/internal/dtos"
	"final-project-backend/internal/helpers"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllSubscriptions(c *gin.Context) {
	var response dtos.GetAllSubscriptionsResponse

	subscriptions, err := h.services.Subscription.GetAll()
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	response = dtos.FormatSubscriptions(subscriptions)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) AddUserSubscription(c *gin.Context) {
	var request dtos.AddUserSubscriptionRequest
	var response dtos.AddUserSubscriptionResponse

	err := c.ShouldBindJSON(&request)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	userSubscription, err := h.services.UserSubscription.AddUserSubscription(request.UserID, request.SubscriptionID)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	response = dtos.AddUserSubscriptionResponse{
		UserID:         userSubscription.UserID,
		SubscriptionID: userSubscription.SubscriptionID,
		DateStarted:    userSubscription.DateStarted,
		DateEnded:      userSubscription.DateEnded,
		RemainingQuota: userSubscription.RemainingQuota,
	}

	helpers.SendSuccessResponse(c, http.StatusCreated, http.StatusText(http.StatusCreated), response)
}
