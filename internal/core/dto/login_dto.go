package dto

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // username or email
	Password   string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
