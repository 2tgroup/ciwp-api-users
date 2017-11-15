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
	Name         string        `json:"name"`
	Email        string        `json:"email"`
	PasswordHash string        `json:"password_hash"`
	UserType     bool          `json:"user_type"`
	UserInfo     UserInfo      `json:"user_info"`
	Meta         interface{}   `json:"meta"`
}

//UserInfo hold billing info
type UserInfo struct {
	ExtendKeyCard  string           `json:"extend_key_card"`
	ListCards      []typeUserCard   `json:"cards"`
	Wallets        []typeUserWallet `json:"wallets"`
	CurrentCard    typeUserCard     `json:"current_card"`
	CurrentWallets typeUserWallet   `json:"current_wallet"`
	CurrentBlance  float32          `json:"current_blance"`
	Address        UserAddress      `json:"address"`
}

/*UserAddress ex: address,phone... */
type UserAddress struct {
	Street   string `json:"street"`
	City     string `json:"city"`
	State    string `json:"state"`
	PostCode string `json:"post_code"`
	Country  string `json:"country"`
}

type typeUserCard struct {
	CardName        string   `json:"card_name"`
	CardLastDigital string   `json:"last_digital"`
	CardExpried     string   `json:"expried_date"`
	CustomerID      []string `json:"customer_ids"`
}

type typeUserWallet struct {
	WalletName string            `json:"wl_name"`
	WalletIDs  map[string]string `json:"wl_ids"`
}

/*UserGeneratePass crypt password*/
func (userBase *UserBase) UserGeneratePass() {
	cryPass, err := bcrypt.GenerateFromPassword([]byte(userBase.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		// TODO: Properly handle error
		log.Fatal(err)
	}
	userBase.PasswordHash = string(cryPass)
}

//UserCheckPass is verifly password for user
func (userBase *UserBase) UserCheckPass(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userBase.PasswordHash), []byte(password))
	return err == nil
}

/*UserAdd Insert user*/
func (userBase *UserBase) UserAdd() error {
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
