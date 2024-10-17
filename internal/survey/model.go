package survey

import "time"

type Survey struct {
	Id        int        `json:"id"`
	UserId    int        `json:"user_id"`
	Name      string     `json:"name"`
	PhoneNo   string     `json:"phone_no"`
	Gender    string     `json:"gender"`
	Dob       *time.Time `json:"dob"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type SurveyInput struct {
	Id        int    `json:"id"`
	UserId    int    `json:"user_id"`
	Name      string `json:"name"`
	PhoneNo   string `json:"phone_no"`
	Gender    string `json:"gender"`
	Dob       string `json:"dob"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
