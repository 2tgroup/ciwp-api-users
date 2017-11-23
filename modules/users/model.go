package users

import (
	"encoding/json"
	"log"
	"time"

	"bitbucket.org/2tgroup/ciwp-api-users/config"
	"bitbucket.org/2tgroup/ciwp-api-users/dbconnects"
	"bitbucket.org/2tgroup/ciwp-api-users/types"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var collection string
var statusUser int

func init() {
	collection = "users"
	statusUser = 1
}

/*UserBase it can be extend */
type UserBase struct {
	ID           bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string        `json:"name,omitempty" bson:"name,omitempty"`
	Email        string        `json:"email" validate:"required,email" bson:"email"`
	Password     string        `json:"password,omitempty" validate:"required" bson:"password,omitempty"`
	PasswordHash string        `json:"password_hash,omitempty" bson:"password_hash,omitempty"`
	UserType     string        `json:"user_type,omitempty" bson:"user_type,omitempty"`
	UserInfo     UserInfo      `json:"user_info,omitempty" bson:"user_info"`
	Status       int           `json:"status,omitempty" bson:"status,omitempty"`
	Meta         interface{}   `json:"meta,omitempty" bson:"meta,omitempty"`
	Create       time.Time     `json:"created,omitempty" bson:"created,omitempty"`
	Updated      time.Time     `json:"updated,omitempty" bson:"updated,omitempty"`
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

func (userBase *UserBase) defaultValueUser() {
	userBase.Status = statusUser
	userBase.UserType = "mem"
	if userBase.UserInfo.Currency == "" {
		userBase.UserInfo.Currency = "usd"
	}
	userBase.Create = time.Now()
	userBase.Updated = time.Now()
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
func (userBase *UserBase) UserCheckPass() bool {
	err := bcrypt.CompareHashAndPassword([]byte(userBase.PasswordHash), []byte(userBase.Password))
	return err == nil
}

/*UserAdd Insert user*/
func (userBase *UserBase) UserAdd() error {
	userBase.UserGeneratePass()
	userBase.defaultValueUser()
	if err := dbconnect.InserToCollection(collection, userBase); err != nil {
		return err
	}
	userBase.UserGetOne(echo.Map{
		"email": userBase.Email,
	})
	return nil
}

/*UserAddAdmin add admin manager*/
func (userBase *UserBase) UserAddAdmin() error {
	userBase.UserType = "admin"
	return userBase.UserAdd()
}

/*UserUpdate Update users*/
func (userBase *UserBase) UserUpdate(_id string) error {

	condition := dbconnect.MongodbToBson(echo.Map{
		"_id": bson.ObjectIdHex(_id),
	})

	if userBase.Password != "" {
		userBase.UserGeneratePass()
	}
	// delete field
	userBase.Password = ""
	userBase.UserType = ""
	userBase.Updated = time.Now()
	dataSet := dbconnect.MongodbToBson(echo.Map{
		"$set": userBase,
	})
	return dbconnect.UpdateOneInCollection(collection, condition, dataSet)
}

/*UserGetOne get single user*/
func (userBase *UserBase) UserGetOne(q interface{}) error {
	/*Convert query to bson*/
	query := dbconnect.MongodbToBson(q)

	userDataRaw, err := dbconnect.GetOneDataInCollection(collection, query)

	if err != nil {
		return err
	}

	byteData, errMar := json.Marshal(userDataRaw)

	if errMar != nil {
		return err
	}
	json.Unmarshal(byteData, &userBase)

	return nil
}

//AuthSignupToken genarate token
func (userBase *UserBase) AuthSignupToken() (string, error) {

	a := types.AuthJwtClaims{}
	a.StandardClaims = jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	}
	a.ID = userBase.ID
	a.Name = userBase.Name
	a.Email = userBase.Email
	a.UserType = userBase.UserType
	a.Info = userBase.UserInfo
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.DataConfig.SecretKey))
	return t, err
}
