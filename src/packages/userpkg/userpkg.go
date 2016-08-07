package userpkg

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"not null;unique"`
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
}

func CreateUser(db *gorm.DB, username string, email string, password string) *User {
	//hash the password eventually
	newUser := User{Username: username, Email: email, Password: password}
	if db.NewRecord(newUser) == true {
		if err := db.Create(&newUser).Error; err != nil {
			return nil
		} else {
			return &newUser
		}
	}
	return nil
}
