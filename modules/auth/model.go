package auth

import (
	"bitbucket.org/2tgroup/ciwp-api-users/modules/users"
)

type authUser struct {
	ID          string      `json:"_id"`
	Email       string      `json:"email"`
	Name        string      `json:"name"`
	UserType    string      `json:"user_type"`
	Avatar      string      `json:"avatar"`
	Status      int         `json:"status"`
	SesstionExp int64       `json:"session_exp"`
	Info        interface{} `json:"info"`
}

type AuthResponse struct {
	Token    string   `json:"token"`
	UserInfo authUser `json:"user"`
}

//AuthSetResponse format data response
func (au *AuthResponse) AuthSetResponse(user users.UserBase) {
	au.UserInfo.ID = user.ID.Hex()
	au.UserInfo.Name = user.Name
	au.UserInfo.Email = user.Email
	au.UserInfo.UserType = user.UserType
	au.UserInfo.Info = user.UserInfo
	au.UserInfo.Status = user.Status
	au.UserInfo.Avatar = user.Avatar
	au.UserInfo.SesstionExp = user.SesstionExp
}
