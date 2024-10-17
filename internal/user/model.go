package user

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Token struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
	ExpiredAt   int64  `json:"expired_at"`
}
