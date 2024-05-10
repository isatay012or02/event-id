package models

type AuthEvent struct {
	AID         string `json:"aid"`
	RedirectURL string `json:"redirectURL"`
	Status      string `json:"status"`
	ClientID    string `json:"clientID"`
}
