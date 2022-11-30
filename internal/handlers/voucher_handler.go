package handlers

import (
	"net/http"

	"final-project-backend/internal/dtos"
	"final-project-backend/internal/helpers"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllVouchers(c *gin.Context) {
	var response dtos.GetAllVouchersResponse

	vouchers, err := h.services.Voucher.GetAll()
	if err != nil {
		helpers.SendErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))

		return
	}

	response = dtos.FormatVouchers(vouchers)

	helpers.SendSuccessResponse(c, http.StatusOK, http.StatusText(http.StatusOK), response)
}
