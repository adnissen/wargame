package userpkg

import (
	"crypto/sha256"
	"io"

	"github.com/adnissen/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null;unique"`
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	Gold     int
}

func CreateUser(db *gorm.DB, username string, email string, password string) *User {
	if len(password) < 7 {
		return nil
	}
	h := sha256.New()
	io.WriteString(h, password)
	p := string(h.Sum(nil))

	//hash the password eventually
	newUser := User{Username: username, Email: email, Password: p}

	if db.NewRecord(newUser) == true {
		if err := db.Create(&newUser).Error; err != nil {
			return nil
		} else {
			return &newUser
		}
	}
	return nil
}

func VerifyUser(db *gorm.DB, username string, password string) *User {
	h := sha256.New()
	io.WriteString(h, password)
	p := string(h.Sum(nil))

	var u *User
	db.Where(&User{Username: username, Password: p}).First(u)
	if u != nil {
		return u
	}
	return nil
}
