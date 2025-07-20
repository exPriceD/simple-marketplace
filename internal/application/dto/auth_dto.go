package dto

type AuthDTO struct {
	AccessToken string  `json:"access_token"`
	User        UserDTO `json:"user"`
}
