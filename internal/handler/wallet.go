package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"newWalletService/internal/model"
	"newWalletService/internal/usecase"
)

type WalletHandler struct {
	walletUsecase *usecase.WalletUsecase
}

func NewWalletHandler(walletUsecase *usecase.WalletUsecase) *WalletHandler {
	return &WalletHandler{
		walletUsecase: walletUsecase,
	}
}

func (h *WalletHandler) Create(c *gin.Context) {
	var wallet model.Wallet

	if err := c.ShouldBindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input format",
			"details": err.Error(),
		})
		return
	}

	id, err := h.walletUsecase.Create(wallet.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Wallet creation failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Wallet registered successfully",
		"wallet_id": id,
	})
}

//
//// Login handles wallet authentication and JWT generation
//func (h *WalletHandler) Login(c *gin.Context) {
//	var login model.WalletLogin
//	if err := c.ShouldBindJSON(&login); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login data"})
//		return
//	}
//
//	// Get wallet from database
//	var wallet model.Wallet
//	wallet, err := h.walletUsecase.Find(login.Walletname)
//
//	if err == sql.ErrNoRows {
//		// Don't specify whether email or password was wrong
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
//		return
//	}
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login process failed"})
//		return
//	}
//
//	// Verify password
//	if !h.walletUsecase.CheckPassword(login.Password, wallet.Password) {
//		// Use same message as above for security
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
//		return
//	}
//
//	// Generate JWT with claims
//	tokenString, err := h.walletUsecase.GenerateJWT(wallet, h.tokenExpiration, h.jwtSecret)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
//		return
//	}
//
//	// Return token with expiration
//	c.JSON(http.StatusOK, gin.H{
//		"token":      tokenString,
//		"expires_in": h.tokenExpiration.Seconds(),
//		"token_type": "Bearer",
//	})
//}
