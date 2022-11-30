package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"final-project-backend/internal/constants"
	"final-project-backend/internal/dtos"
	errn "final-project-backend/internal/errors"
	"final-project-backend/internal/helpers"
	"final-project-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) GetAllGifts(c *gin.Context) {
	var response dtos.GetAllGiftsResponse

	gifts, err := h.services.Gift.GetAll()
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	response = dtos.FormatGifts(gifts)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) GetAllUserGifts(c *gin.Context) {
	var response dtos.GetAllUserGiftsResponse

	userGifts, err := h.services.Gift.GetAllUserGifts()
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	response = dtos.FormatUserGifts(userGifts)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) UpdateUserGiftStatus(c *gin.Context) {
	var request dtos.UpdateProcessedInvoiceRequest
	var response dtos.UpdateUserGiftStatusResponse

	err := c.ShouldBindJSON(&request)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	id := c.Param("id")
	idSplitted := strings.Split(id, "-")
	if len(idSplitted) != 4 {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	userID, err1 := strconv.ParseInt(idSplitted[0], 10, 64)
	giftID, err2 := strconv.ParseInt(idSplitted[1], 10, 64)
	month, err3 := strconv.Atoi(idSplitted[2])
	year, err4 := strconv.Atoi(idSplitted[3])

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
	}

	var newStatus constants.GiftStatus
	if *request.IsSuccess {
		newStatus = constants.COMPLETED
	} else {
		newStatus = constants.CANCELLED
	}

	newUserGifts := &models.UserGift{
		UserID: userID,
		GiftID: giftID,
		Month:  month,
		Year:   year,
		Status: newStatus,
	}

	updatedUserGift, err := h.services.Gift.UpdateUserGiftStatus(newUserGifts)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.SendErrorResponse(c, http.StatusNotFound, errn.ErrNoInvoicesFound.Error())
			return
		}

		if errors.Is(err, errn.ErrInvalidStatusUpdate) {
			helpers.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	response = dtos.FormatUserGift(updatedUserGift)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) UpdateGiftStock(c *gin.Context) {
	var request dtos.UpdateGiftStockRequest
	var response dtos.UpdateGiftStockResponse

	err := c.ShouldBindJSON(&request)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	id := c.Param("id")
	parsedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	updatedGift, err := h.services.Gift.UpdateGiftStock(parsedID, request.Stock)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.SendErrorResponse(c, http.StatusNotFound, errn.ErrNoInvoicesFound.Error())
			return
		}

		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	response = dtos.FormatGift(updatedGift)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}
