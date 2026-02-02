package handler

import (
	"meye-core/internal/application/user"
	dto "meye-core/internal/infrastructure/api/handler/dto/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	createUserUseCase user.CreateUserUseCase
	loginUseCase      user.LoginUseCase
}

func NewUserHandler(createUserUC user.CreateUserUseCase, loginUseCase user.LoginUseCase) *UserHandler {
	return &UserHandler{
		createUserUseCase: createUserUC,
		loginUseCase:      loginUseCase,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var reqBody dto.CreateUserInputBody

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := user.CreateUserInput{
		Username: reqBody.Username,
		Password: reqBody.Password,
		Role:     reqBody.Role,
	}

	user, err := h.createUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		respondMappedError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.MapUserOutput(user))
}

func (h *UserHandler) Login(c *gin.Context) {
	var reqBody dto.LoginInputBody

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := user.LoginInput{
		Username: reqBody.Username,
		Password: reqBody.Password,
	}

	token, err := h.loginUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		respondMappedError(c, err)
		return
	}

	output := dto.MapLoginOutput(token)

	c.JSON(http.StatusOK, output)
}
