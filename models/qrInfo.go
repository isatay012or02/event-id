package models

type QRInfo struct {
	ClientID    string `json:"clientID"`
	Scope       string `json:"scope"`
	RedirectURL string `json:"redirectURL"`
	State       string `json:"state"`
	AID         string `json:"aid"`
}
