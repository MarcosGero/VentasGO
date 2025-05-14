package api

import (
	"net/http"
	"parte3/internal/sale"

	"github.com/gin-gonic/gin"
)

type saleHandler struct {
	svc *sale.Service
}

func (h *saleHandler) handleCreate(ctx *gin.Context) {
	var req struct {
		UserID string  `json:"user_id"`
		Amount float64 `json:"amount"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s, err := h.svc.Create(req.UserID, req.Amount)
	if err != nil {
		switch err {
		case sale.ErrBadAmount, sale.ErrUserNotFound:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusCreated, s)
}

func (h *saleHandler) handlePatch(ctx *gin.Context) {
	id := ctx.Param("id")
	var req struct {
		Status sale.Status `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s, err := h.svc.UpdateStatus(id, req.Status)
	if err != nil {
		switch err {
		case sale.ErrNotFound:
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case sale.ErrBadStatus:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case sale.ErrBadTrans:
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, s)
}

func (h *saleHandler) handleSearch(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	var statusPtr *sale.Status
	if raw := ctx.Query("status"); raw != "" {
		st := sale.Status(raw)
		statusPtr = &st
	}
	res, err := h.svc.Search(userID, statusPtr)
	if err != nil {
		if err == sale.ErrBadStatus {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}
