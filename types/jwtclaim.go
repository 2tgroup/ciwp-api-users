package types

import (
	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

/*AuthJwtClaims customize Auth Claims*/
type AuthJwtClaims struct {
	ID       bson.ObjectId          `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string                 `json:"name"`
	Avatar   string                 `json:"avatar"`
	Email    string                 `json:"email"`
	UserType string                 `json:"user_type,omitempty" bson:"user_type,omitempty"`
	Info     interface{}            `json:"info"`
	Extends  map[string]interface{} `json:"meta"`
	jwt.StandardClaims
}
