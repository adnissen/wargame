package invitecode

import (
	"math/rand"
	"time"

	"github.com/adnissen/gorm"
)

type InviteCode struct {
	gorm.Model
	Code     string `gorm:"not null;unique"`
	UserId   uint
	Redeemed bool
}

func CreateCode(db *gorm.DB) *InviteCode {
	code := InviteCode{Code: RandStringBytesMaskImprSrc(24)}
	if err := db.Create(&code).Error; err != nil {
		return nil
	} else {
		return &code
	}
}

func VerifyCode(db *gorm.DB, code string) *InviteCode {
	ic := InviteCode{}
	if err := db.Where(&InviteCode{Code: code}).First(&ic).Error; err != nil && !ic.Redeemed {
		return &ic
	}
	return &ic
}

func (ic *InviteCode) Claim(db *gorm.DB, id uint) bool {
	if !ic.Redeemed && ic.UserId == 0 {
		ic.Redeemed = true
		ic.UserId = id
		db.Save(ic)
		return true
	}
	return false
}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
