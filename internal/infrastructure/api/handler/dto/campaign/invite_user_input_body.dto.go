package campaign

type InviteUserInputBody struct {
	UserID string `json:"user_id" binding:"required"`
}