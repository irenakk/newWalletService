package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"walletService/internal/dto"
	"walletService/internal/usecase"
)

type WalletHandler struct {
	walletUsecase *usecase.WalletUsecase
}

func NewWalletHandler(walletUsecase *usecase.WalletUsecase) *WalletHandler {
	return &WalletHandler{
		walletUsecase: walletUsecase,
	}
}

func (h *WalletHandler) Add(c *gin.Context) {
	var addRequest dto.Add

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&addRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input format",
			"details": err.Error(),
		})
		return
	}

	account, err := h.walletUsecase.Add(username.(string), addRequest.Currency, addRequest.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Amount adding failed"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Amount successfully added",
		"account": account,
	})
}

func (h *WalletHandler) Transfer(c *gin.Context) {
	var addRequest dto.Transfer

	usernameFrom, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&addRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input format",
			"details": err.Error(),
		})
		return
	}

	account, err := h.walletUsecase.Transfer(usernameFrom.(string), addRequest.Username, addRequest.Currency, addRequest.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Amount adding failed"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Amount successfully added",
		"account": account,
	})
}
