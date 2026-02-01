package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func respondMappedError(c *gin.Context, err error) {
	logrus.WithContext(c.Request.Context()).Error(err.Error())
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: "Internal Server Error",
		Code:  "INTERNAL_SERVER_ERROR",
	})
}
