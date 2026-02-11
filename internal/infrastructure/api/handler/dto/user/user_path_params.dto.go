package dto

type UserPathParams struct {
	UserID string `uri:"userID" binding:"required"`
}
