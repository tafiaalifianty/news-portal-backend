package dtos

type RegisterRequestDTO struct {
	Email        string  `json:"email" binding:"required"`
	Password     string  `json:"password" binding:"required"`
	Fullname     string  `json:"fullname" binding:"required"`
	Address      string  `json:"address" binding:"required"`
	ReferralCode *string `json:"referral_code"`
}

type RegisterResponseDTO struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	Fullname     string `json:"fullname"`
	Address      string `json:"address"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginRequestDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponseDTO struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type JwtData struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Token string `json:"token"`
}

type RefreshTokenData struct {
	ID int64 `json:"id"`
}

type RefreshTokenContext struct {
	ID           int64  `json:"id"`
	RefreshToken string `json:"refresh_token"`
}
