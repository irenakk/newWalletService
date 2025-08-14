package handler

import (
	//"github.com/gin-gonic/gin"
	//"net/http"
	//"newWalletService/internal/model"
	"newWalletService/internal/usecase"
)

type AccountHandler struct {
	accountUsecase *usecase.AccountUsecase
}

func NewAccountHandler(accountUsecase *usecase.AccountUsecase) *AccountHandler {
	return &AccountHandler{
		accountUsecase: accountUsecase,
	}
}

//func (h *AccountHandler) Create(c *gin.Context) {
//	var account model.Account
//
//	if err := c.ShouldBindJSON(&account); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"error":   "Invalid input format",
//			"details": err.Error(),
//		})
//		return
//	}
//
//	id, err := h.accountUsecase.Create(account.ID)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Account creation failed"})
//		return
//	}
//
//	c.JSON(http.StatusCreated, gin.H{
//		"message":   "Account registered successfully",
//		"account_id": id,
//	})
//}
