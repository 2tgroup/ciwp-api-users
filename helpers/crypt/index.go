package helperCtypt

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

//GenerateCrypt generate hash string
func GenerateCrypt(stringText string) string {
	cryText, err := bcrypt.GenerateFromPassword([]byte(stringText), bcrypt.DefaultCost)
	if err != nil {
		// TODO: Properly handle error
		log.Fatal(err)
	}
	return string(cryText)
}

//CheckCrypt is verifly hash
func CheckCrypt(stringCheck, stringCrypt string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(stringCrypt), []byte(stringCheck))
	return err == nil
}
