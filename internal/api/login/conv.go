package login

func convToResponse(token string) Output {
	return Output{Token: token}
}
