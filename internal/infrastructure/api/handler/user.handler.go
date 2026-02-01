package handler

import (
	"meye-core/internal/application/user/createuser"
	"meye-core/internal/infrastructure/api/handler/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	createUserUseCase *createuser.UseCase
}

func NewUserHandler(createUserUC *createuser.UseCase) *UserHandler {
	return &UserHandler{
		createUserUseCase: createUserUC,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var reqBody dto.CreateUserInput

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := createuser.Input{
		Username: reqBody.Username,
		Password: reqBody.Password,
		Role:     reqBody.Role,
	}

	result, err := h.createUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		respondMappedError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.MapUserOutput(result))
}
