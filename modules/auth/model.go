package auth

import "github.com/dgrijalva/jwt-go"

/*Customize Auth Claims*/

type AuthJwtClaims struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"isAdmin"`
	Info    map[string]interface{}
	Extends map[string]interface{} `json:"extends"`
	jwt.StandardClaims
}
