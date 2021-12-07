package login_dto

import "time"

type LoginRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=8"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token" binding:"required"`
}

type LoginResponse struct {
	ID                    uint      `json:"id"`
	Name                  string    `json:"name"`
	Email                 string    `json:"email"`
	AccessToken           string    `json:"access_token"`
	AccesssTokenExpiredAt time.Time `json:"access_token_expired_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiredAt time.Time `json:"refresh_token_expired_at"`
}

type RefreshAccessTokenResponse struct {
	ID                    uint      `json:"id"`
	Name                  string    `json:"name"`
	Email                 string    `json:"email"`
	AccessToken           string    `json:"access_token"`
	AccesssTokenExpiredAt time.Time `json:"access_token_expired_at"`
	Message               string    `json:"message"`
}
