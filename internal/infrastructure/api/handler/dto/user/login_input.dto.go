package dto

type LoginInput struct {
	Username string `json:"username" binding:"required,alphanum,min=3,max=100"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}
