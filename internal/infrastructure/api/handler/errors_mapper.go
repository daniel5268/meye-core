package handler

import (
	"errors"
	applicationcampaign "meye-core/internal/application/campaign"
	applicationuser "meye-core/internal/application/user"
	domaincampaign "meye-core/internal/domain/campaign"
	domainuser "meye-core/internal/domain/user"
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
	case errors.Is(err, applicationuser.ErrUserNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "User not found",
			Code:  applicationuser.ErrUserNotFound.Error(),
		})
	case errors.Is(err, applicationcampaign.ErrCampaignNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Campaign not found",
			Code:  applicationcampaign.ErrCampaignNotFound.Error(),
		})
	case errors.Is(err, domainuser.ErrUserNotPlayer):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "User is not a player",
			Code:  domainuser.ErrUserNotPlayer.Error(),
		})
	case errors.Is(err, domaincampaign.ErrUserNotInvited):
		c.JSON(http.StatusNotAcceptable, ErrorResponse{
			Error: "User should be invited to create PJs in a campaign",
			Code:  domaincampaign.ErrUserNotInvited.Error(),
		})
	case errors.Is(err, domaincampaign.ErrPJsNotInCampaign):
		c.JSON(http.StatusNotAcceptable, ErrorResponse{
			Error: "All PJs should belong to the campaign",
			Code:  domaincampaign.ErrPJsNotInCampaign.Error(),
		})
	case errors.Is(err, domaincampaign.ErrInsufficientXP):
		c.JSON(http.StatusNotAcceptable, ErrorResponse{
			Error: "There is not enough XP to perform the action",
			Code:  domaincampaign.ErrInsufficientXP.Error(),
		})
	case errors.Is(err, domaincampaign.ErrCannotReduceStats):
		c.JSON(http.StatusNotAcceptable, ErrorResponse{
			Error: "PJ stats can't be reduced",
			Code:  domaincampaign.ErrCannotReduceStats.Error(),
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal Server Error",
			Code:  "INTERNAL_SERVER_ERROR",
		})
	}
}
