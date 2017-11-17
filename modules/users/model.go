package users

import (
	"encoding/json"
	"log"

	"bitbucket.org/2tgroup/ciwp-api-users/dbconnects"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var collection string

func init() {
	collection = "users"
}

/*UserBase it can be extend */
type UserBase struct {
	ID           bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string        `json:"name" validate:"required" bson:"name"`
	Email        string        `json:"email" validate:"required,email" bson:"email"`
	Password     string        `json:"password,omitempty" validate:"required" bson:"password,omitempty"`
	PasswordHash string        `json:"password_hash,omitempty" bson:"password_hash,omitempty"`
	UserType     string        `json:"user_type,omitempty" bson:"user_type,omitempty"`
	UserInfo     UserInfo      `json:"user_info,omitempty" bson:"user_info"`
	Meta         interface{}   `json:"meta,omitempty" bson:"meta,omitempty"`
}

//UserInfo hold billing info
type UserInfo struct {
	ExtendKeyCard  string           `json:"extend_key_card,omitempty" bson:"extend_key_card"`
	ListCards      []typeUserCard   `json:"cards,omitempty" bson:"cards"`
	Wallets        []typeUserWallet `json:"wallets,omitempty" bson:"wallets"`
	CurrentCard    typeUserCard     `json:"current_card,omitempty" bson:"current_card"`
	CurrentWallets typeUserWallet   `json:"current_wallet,omitempty" bson:"current_wallet"`
	CurrentBlance  float32          `json:"current_blance,omitempty" bson:"current_blance"`
	Address        UserAddress      `json:"address,omitempty" bson:"address"`
	Currency       string           `json:"currency,omitempty"`
}

/*UserAddress ex: address,phone... */
type UserAddress struct {
	Street   string `json:"street,omitempty" bson:"street"`
	City     string `json:"city,omitempty" bson:"city"`
	State    string `json:"state,omitempty" bson:"state"`
	PostCode string `json:"post_code,omitempty" bson:"post_code"`
	Country  string `json:"country,omitempty" bson:"country"`
}

type typeUserCard struct {
	CardName        string   `json:"card_name,omitempty" bson:"card_name"`
	CardLastDigital string   `json:"last_digital,omitempty" bson:"last_digital"`
	CardExpried     string   `json:"expried_date,omitempty" bson:"expried_date"`
	CustomerID      []string `json:"customer_ids,omitempty" bson:"customer_ids"`
}

type typeUserWallet struct {
	WalletName string            `json:"wl_name,omitempty" bson:"wl_name"`
	WalletIDs  map[string]string `json:"wl_ids,omitempty" bson:"wl_ids"`
}

/*UserGeneratePass crypt password*/
func (userBase *UserBase) UserGeneratePass() {
	cryPass, err := bcrypt.GenerateFromPassword([]byte(userBase.Password), bcrypt.DefaultCost)
	if err != nil {
		// TODO: Properly handle error
		log.Fatal(err)
	}
	userBase.PasswordHash = string(cryPass)
	userBase.Password = ""
}

//UserCheckPass is verifly password for user
func (userBase *UserBase) UserCheckPass(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userBase.PasswordHash), []byte(password))
	return err == nil
}

/*UserAdd Insert user*/
func (userBase *UserBase) UserAdd() error {
	userBase.UserType = "normal"
	userBase.UserGeneratePass()
	if userBase.UserInfo.Currency == "" {
		userBase.UserInfo.Currency = "usd"
	}
	return dbconnect.InserToCollection(collection, userBase)
}

/*UserUpdate Update users*/
func (userBase *UserBase) UserUpdate() error {
	return nil
}

/*UserGetOne get single user*/
func (userBase *UserBase) UserGetOne(q interface{}) (userData *UserBase, err error) {
	/*Convert query to bson*/
	query := dbconnect.MongodbToBson(q)

	userDataRaw, err := dbconnect.GetOneDataInCollection(collection, query)

	if err != nil {
		return nil, err
	}

	byteData, errMar := json.Marshal(userDataRaw)

	if errMar != nil {
		return nil, err
	}

	json.Unmarshal(byteData, &userData)

	return userData, err
}
