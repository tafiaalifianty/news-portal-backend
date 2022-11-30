package handlers

import (
	"errors"
	"net/http"
	"strings"

	"final-project-backend/internal/dtos"
	errn "final-project-backend/internal/errors"
	"final-project-backend/internal/helpers"
	"final-project-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) GetAllInvoices(c *gin.Context) {
	var response dtos.GetAllInvoicesResponse

	invoices, err := h.services.Invoice.GetAll()
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	response = dtos.FormatInvoices(invoices)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) GetWaitingInvoiceByCode(c *gin.Context) {
	var response dtos.InvoiceResponse

	code := c.Param("code")

	invoice, err := h.services.Invoice.GetByCode(strings.ToUpper(code))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.SendErrorResponse(c, http.StatusNotFound, errn.ErrNoInvoicesFound.Error())
			return
		}
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	if invoice.Status != models.WAITING {
		helpers.SendErrorResponse(c, http.StatusBadRequest, errn.ErrInvoiceNotAwaitingPayment.Error())
		return
	}

	response = *dtos.FormatInvoice(invoice)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) GetWaitingInvoiceByCodeProtected(c *gin.Context) {
	userContext, ok := c.Get("user")
	if !ok {
		helpers.SendErrorResponse(
			c,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	var response dtos.InvoiceResponse

	code := c.Param("code")

	invoice, err := h.services.Invoice.GetUserInvoiceByCode(strings.ToUpper(code), userContext.(dtos.JwtData).ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.SendErrorResponse(c, http.StatusNotFound, errn.ErrNoInvoicesFound.Error())
			return
		}

		if errors.Is(err, errn.ErrNotAuthorized) {
			helpers.SendErrorResponse(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	if invoice.Status != models.WAITING {
		helpers.SendErrorResponse(c, http.StatusBadRequest, errn.ErrInvoiceNotAwaitingPayment.Error())
		return
	}

	response = *dtos.FormatInvoice(invoice)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) CreateInvoice(c *gin.Context) {
	var request dtos.CreateInvoiceRequest
	var response dtos.CreateInvoiceResponse

	err := c.ShouldBindJSON(&request)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	userContext, ok := c.Get("user")
	if !ok {
		helpers.SendErrorResponse(
			c,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
		)
		return
	}

	user := userContext.(dtos.JwtData)
	createdInvoice, err := h.services.Invoice.Create(&models.Invoice{
		UserID:         user.ID,
		OriginalPrice:  request.OriginalPrice,
		Total:          request.OriginalPrice,
		Status:         models.WAITING,
		SubscriptionID: int64(request.SubscriptionID),
		VoucherCode:    request.VoucherCode,
	})
	if err != nil {
		if errors.Is(err, errn.ErrInvalidVoucher) {
			helpers.SendErrorResponse(
				c,
				http.StatusBadRequest,
				err.Error(),
			)
			return
		}

		if errors.Is(err, errn.ErrVoucherExpired) {
			helpers.SendErrorResponse(
				c,
				http.StatusBadRequest,
				err.Error(),
			)
			return
		}

		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	response = *dtos.FormatInvoice(createdInvoice)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) UpdateWaitingInvoiceToProcessed(c *gin.Context) {
	var response dtos.InvoiceResponse

	code := c.Param("code")

	updatedInvoice, _, _, err := h.services.Invoice.UpdateStatus(strings.ToUpper(code), models.PROCESSED)
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

	response = *dtos.FormatInvoice(updatedInvoice)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}

func (h *Handler) UpdateProcessedInvoice(c *gin.Context) {
	var request dtos.UpdateProcessedInvoiceRequest
	var response dtos.UpdateProcessedInvoiceResponse

	err := c.ShouldBindJSON(&request)
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	code := c.Param("code")

	var newStatus models.InvoiceStatus
	if *request.IsSuccess {
		newStatus = models.COMPLETED
	} else {
		newStatus = models.REJECTED
	}

	updatedInvoice, giftsSent, vouchersSent, err := h.services.Invoice.UpdateStatus(strings.ToUpper(code), newStatus)
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

	response = dtos.UpdateProcessedInvoiceResponse{
		InvoiceResponse: *dtos.FormatInvoice(updatedInvoice),
		GiftsSent:       dtos.FormatGifts(giftsSent),
		VouchersSent:    dtos.FormatVouchers(vouchersSent),
	}

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}
