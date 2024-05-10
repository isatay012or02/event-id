package models

type RedirectURLResponse struct {
	CodeMessage
	RedirectURL string
}

func RedirectURLSuccessResponse(redirectURL string) *RedirectURLResponse {

	return &RedirectURLResponse{
		CodeMessage: CodeMessage{
			Code:    0,
			Message: "OK",
		},
		RedirectURL: redirectURL,
	}

}
