package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"parte3/internal/sale"
)

type SaleHandler struct {
	service sale.Service
}

func NewSaleHandler(s sale.Service) *SaleHandler {
	return &SaleHandler{service: s}
}

type createSaleRequest struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

func (h *SaleHandler) Create(c *gin.Context) {
	var req createSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	s, err := h.service.Create(req.UserID, req.Amount)
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case sale.ErrUserNotFound:
			status = http.StatusNotFound
		case sale.ErrBadAmount:
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, s)
}
