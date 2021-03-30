package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mehrdod/todo/domain"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var user domain.User

	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.CreateUser(user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (h *Handler) signIn(c *gin.Context) {

}
