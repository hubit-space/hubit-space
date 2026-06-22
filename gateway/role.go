package model

import "github.com/golang-jwt/jwt/v5"

type AccessToken struct {
	UserId       string `json:"user_id"`
	UserCode     string `json:"user_code"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	SalesCode    string `json:"sales_code"`
	PositionCode string `json:"position_code"`
	RoleCode     string `json:"role_code"`
	NIK          string `json:"nik"`
	jwt.RegisteredClaims
}
