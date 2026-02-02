package handler

import (
	"errors"
	applicationuser "meye-core/internal/application/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func respondMappedError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, applicationuser.ErrUsernameAlreadyExists):
		c.JSON(http.StatusConflict, ErrorResponse{
			Error: "Username already exists",
			Code:  applicationuser.ErrUsernameAlreadyExists.Error(),
		})
	case errors.Is(err, applicationuser.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Invalid credentials",
			Code:  applicationuser.ErrInvalidCredentials.Error(),
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal Server Error",
			Code:  "INTERNAL_SERVER_ERROR",
		})
	}
}
