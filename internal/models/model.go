package models

type Token struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
	ExpiredAt   int64  `json:"expired_at"`
}

type Response struct {
	Data interface{} `json:"data"`
}
