package users

import (
	"encoding/json"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/eefret/gravatar"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"

	"bitbucket.org/2tgroup/ciwp-api-users/config"
	"bitbucket.org/2tgroup/ciwp-api-users/dbconnects"
	"bitbucket.org/2tgroup/ciwp-api-users/types"
)

var collection string
var statusUser int
var gAvatar *gravatar.Gravatar
var err error

func init() {
	collection = "users"
	statusUser = 1
	gAvatar, err = gravatar.New()
}

/*UserBase it can be extend */
type UserBase struct {
	ID           bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string        `json:"name,omitempty" bson:"name,omitempty"`
	Email        string        `json:"email,omitempty" validate:"required,email" bson:"email,omitempty"`
	Password     string        `json:"password,omitempty" validate:"required" bson:"password,omitempty"`
	PasswordHash string        `json:"password_hash,omitempty" bson:"password_hash,omitempty"`
	UserType     string        `json:"user_type,omitempty" bson:"user_type,omitempty"`
	Avatar       string        `json:"avatar,omitempty" bson:"avatar,omitempty"`
	UserInfo     UserInfo      `json:"info,omitempty" bson:"info,omitempty"`
	Status       int           `json:"status,omitempty" bson:"status,omitempty"`
	Meta         interface{}   `json:"meta,omitempty" bson:"meta,omitempty"`
	Create       time.Time     `json:"created,omitempty" bson:"created,omitempty"`
	Updated      time.Time     `json:"updated,omitempty" bson:"updated,omitempty"`
	SesstionExp  int64         `json:"session_exp,omitempty" bson:"-"`
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
	Phone    string `json:"phone,omitempty" bson:"phone"`
}

type typeUserCard struct {
	CardID          string `json:"card_id,omitempty" bson:"card_id"`
	CardName        string `json:"card_name,omitempty" bson:"card_name"`
	CardLastDigital string `json:"last_digital,omitempty" bson:"last_digital"`
	CardExpried     string `json:"expried_date,omitempty" bson:"expried_date"`
	CustomerID      string `json:"customer_id,omitempty" bson:"customer_id"`
	CardToken       string `json:"card_token,omitempty" bson:"-"`
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
	userBase.UserInfo.Address.Country = "US"
	userBase.UserInfo.CurrentBlance = 0
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

//UserCheckEmailExits true/false
func (userBase *UserBase) UserCheckEmailExits(_id string) bool {
	condition := dbconnect.MongodbToBson(echo.Map{
		"_id": echo.Map{
			"$ne": bson.ObjectIdHex(_id),
		},
		"email": userBase.Email,
	})

	return dbconnect.CountRowsInCollection(collection, condition) > 0
}

/*UserAdd Insert user*/
func (userBase *UserBase) UserAdd() error {
	userBase.UserGeneratePass()
	userBase.defaultValueUser()
	userBase.Avatar = gAvatar.URLParse(userBase.Email)
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
	//updated datetime
	userBase.Updated = time.Now()
	dataSet := dbconnect.MongodbToBson(echo.Map{
		"$set": dbconnect.MongodbToBson(userBase),
	})
	return dbconnect.UpdateOneInCollection(collection, condition, dataSet)
}

/*UserSystemUpdate open Update users*/
func (userBase *UserBase) UserSystemUpdate(_id string) error {

	condition := dbconnect.MongodbToBson(echo.Map{
		"_id": bson.ObjectIdHex(_id),
	})

	if userBase.Password != "" {
		userBase.UserGeneratePass()
	}
	// delete field
	userBase.Password = ""
	//updated datetime
	userBase.Updated = time.Now()
	dataSet := dbconnect.MongodbToBson(echo.Map{
		"$set": dbconnect.MongodbToBson(userBase),
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
	a.Avatar = userBase.Avatar
	a.Name = userBase.Name
	a.Email = userBase.Email
	a.UserType = userBase.UserType
	a.Info = userBase.UserInfo
	userBase.SesstionExp = a.ExpiresAt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.DataConfig.SecretKey))
	return t, err
}

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

type UserStructResponse struct {
	Token    string   `json:"token"`
	UserInfo authUser `json:"user"`
}

//UserResponse response to client struct user
func (userBase *UserBase) UserResponse() *UserStructResponse {
	resAuth := new(UserStructResponse)
	token, _ := userBase.AuthSignupToken()
	resAuth.UserInfo.ID = userBase.ID.Hex()
	resAuth.UserInfo.Name = userBase.Name
	resAuth.UserInfo.Email = userBase.Email
	resAuth.UserInfo.UserType = userBase.UserType
	resAuth.UserInfo.Info = userBase.UserInfo
	resAuth.UserInfo.Status = userBase.Status
	resAuth.UserInfo.Avatar = userBase.Avatar
	resAuth.UserInfo.SesstionExp = userBase.SesstionExp
	return &UserStructResponse{
		Token:    token,
		UserInfo: resAuth.UserInfo,
	}
}
