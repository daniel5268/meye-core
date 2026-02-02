package user

type JWTService interface {
	GenerateSignedToken(user *User) (string, error)
	ValidateToken(tokenString string) (string, error)
}
