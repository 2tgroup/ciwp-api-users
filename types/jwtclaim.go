package types

import (
	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

/*Customize Auth Claims*/

type AuthJwtClaims struct {
	ID       bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string        `json:"name"`
	Email    string        `json:"email"`
	UserType string        `json:"user_type,omitempty" bson:"user_type,omitempty"`
	Info     interface{}
	Extends  map[string]interface{} `json:"meta"`
	jwt.StandardClaims
}
