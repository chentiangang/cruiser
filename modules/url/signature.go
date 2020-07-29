package url

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func (u UrlConfig) Sig() string {
	hash := sha256.New()
	_, err := hash.Write([]byte(u.Request.Data))
	if err != nil {
		panic(err)
	}
	bs := hash.Sum(nil)
	digest := hex.EncodeToString(bs)
	millis := time.Now().Add(300 * time.Second).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"accountId":   u.AccountID(),
		"secretKeyId": u.SecretKeyId(),
		"method":      u.apiMethod(),
		"exp":         millis,
		"digest":      digest,
	})

	signedString, err := token.SignedString([]byte(u.SecretKey()))
	if err != nil {
		panic(err)
	}
	return signedString
}

func (u UrlConfig) AccountID() string {
	if u.Request.AccountId == "" {
		panic("must set accountid")
	}
	return u.Request.AccountId
}

func (u UrlConfig) SecretKeyId() string {
	if u.Request.SecretKeyId == "" {
		panic("secret key id is not exits")
	}

	return u.Request.SecretKeyId
}

func (u UrlConfig) SecretKey() string {
	if u.Request.SecretKey == "" {
		panic("secret key is not exits")
	}
	return u.Request.SecretKey
}

func (u UrlConfig) apiMethod() string {
	if u.Request.ApiMethod == "" {
		panic("please set api_mehtod")
	}

	return u.Request.ApiMethod
}
