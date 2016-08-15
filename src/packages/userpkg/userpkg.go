package userpkg

import (
	"fmt"

	"golang.org/x/crypto/scrypt"

	"github.com/adnissen/gorm"

	"log"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null;unique"`
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	Gold     int
}

const (
	PW_SALT_BYTES = 32
	PW_HASH_BYTES = 64
)

func HashPass(password string) (string, error) {
	salt := []byte("superlksazjdfalsjdfe23password")

	hash, err := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, PW_HASH_BYTES)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash), err
}

func CreateUser(db *gorm.DB, username string, email string, password string) *User {
	if len(password) < 7 {
		return nil
	}
	hashedPass, _ := HashPass(password)
	//hash the password eventually
	newUser := User{Username: username, Email: email, Password: hashedPass}

	if db.NewRecord(newUser) == true {
		if err := db.Create(&newUser).Error; err != nil {
			return nil
		} else {
			fmt.Println(hashedPass)
			return &newUser
		}
	}
	return nil
}

func VerifyUser(db *gorm.DB, username string, password string) *User {
	hashedPass, _ := HashPass(password)
	fmt.Println(hashedPass)
	var user User
	if db.First(&user, "username = ? AND password = ?", username, hashedPass).RecordNotFound() {
		return nil
	}
	return &user
}
