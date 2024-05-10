package models

type UserAuthStatusResponse struct {
	CodeMessage
	UserAuthStatus
}

type UserAuthStatus struct {
	RedirectURL string `json:"redirectURL"`
	Status      string `json:"status"`
}

func UserAuthStatusSuccessResponse(userStatus UserAuthStatus) *UserAuthStatusResponse {
	return &UserAuthStatusResponse{
		CodeMessage:    CodeMessage{Code: 0, Message: "OK"},
		UserAuthStatus: userStatus,
	}
}
