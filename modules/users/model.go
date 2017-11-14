package users

import "gopkg.in/mgo.v2/bson"

/*UserBase it can be extend */
type UserBase struct {
	ID      bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string        `json:"name"`
	Email   string        `json:"email"`
	IsAdmin bool          `json:"isAdmin"`
}

/*UserInfo ex: address,phone... */
type UserInfo struct {
	Street   string `json:"street"`
	City     string `json:"city"`
	State    string `json:"state"`
	PostCode string `json:"post_code"`
	Country  string `json:"country"`
}
