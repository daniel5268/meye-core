package dto

type LoginOutputBody struct {
	Token string `json:"token"`
}

func MapLoginOutput(token string) LoginOutputBody {
	return LoginOutputBody{
		Token: token,
	}
}
