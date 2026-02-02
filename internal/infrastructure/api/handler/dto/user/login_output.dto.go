package dto

type LoginOutput struct {
	Token string `json:"token"`
}

func MapLoginOutput(token string) LoginOutput {
	return LoginOutput{
		Token: token,
	}
}
