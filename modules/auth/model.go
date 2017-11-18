package auth

import (
	"time"

	"bitbucket.org/2tgroup/ciwp-api-users/config"
	"bitbucket.org/2tgroup/ciwp-api-users/modules/users"
	"github.com/dgrijalva/jwt-go"
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

//AuthSignupToken
func (a *AuthJwtClaims) AuthSignupToken(user *users.UserBase) (string, error) {
	a.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	}
	a.ID = user.ID
	a.Name = user.Name
	a.Email = user.Email
	a.UserType = user.UserType
	a.Info = user.UserInfo
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.DataConfig.SecretKey))
	return t, err
}
