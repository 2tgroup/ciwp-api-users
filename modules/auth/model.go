package auth

import "github.com/dgrijalva/jwt-go"
import "gopkg.in/mgo.v2/bson"

/*Customize Auth Claims*/

type AuthJwtClaims struct {
	ID      bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string        `json:"name"`
	Email   string        `json:"email"`
	IsAdmin bool          `json:"isAdmin"`
	Info    map[string]interface{}
	Extends map[string]interface{} `json:"extends"`
	jwt.StandardClaims
}
