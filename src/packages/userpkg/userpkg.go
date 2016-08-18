package userpkg

import (
	"encoding/hex"

	"golang.org/x/crypto/scrypt"

	"github.com/adnissen/gorm"
	_ "github.com/adnissen/gorm/dialects/postgres"

	"github.com/adnissen/wargame/src/packages/army"

	"log"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null;unique"`
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	Gold     int
	Armies   []army.Army
}

const (
	PW_SALT_BYTES = 32
	PW_HASH_BYTES = 64
)

func (u *User) AddArmy(db *gorm.DB, a army.Army) {
	a.UserId = u.ID
	newArmy := army.CreateArmy(db, a)
	if newArmy != nil {
		u.Armies = append(u.Armies, *newArmy)
		db.Save(&u)
	}
}

func (u *User) GetArmy(db *gorm.DB) army.Army {
	var armies []army.Army
	db.Model(u).Related(&armies)
	return armies[0]
}

func HashPass(password string) (string, error) {
	salt := []byte("superlksazjdfalsjdfe23password")

	hash, err := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, PW_HASH_BYTES)
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(hash), err
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
			return &newUser
		}
	}
	return nil
}

func VerifyUser(db *gorm.DB, username string, password string) *User {
	hashedPass, _ := HashPass(password)
	u := new(User)
	db.Where(&User{Username: username, Password: hashedPass}).First(u)
	return u
}
