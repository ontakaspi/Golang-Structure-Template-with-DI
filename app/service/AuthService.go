package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang-structure-template-with-di/app/helper"
	"golang-structure-template-with-di/config"
	"io/ioutil"
	"os"
	"strings"
)

// AuthService JWTService jwt @Service
//
//	type AuthService interface {
//		VerifyJWTRSA(token string) (bool, *jwt.Token, error)
//		UserHasRoles(c *gin.Context, roles ...string) (hasPermission bool)
//		UserHasPermissions(c *gin.Context, Permissions ...string) (hasPermission bool)
//	}
type AuthService struct {
	AuthService *AuthService `di.inject:"AuthService"`
}

// VerifyJWTRSA Verify a JWT token using an RSA public key
func (s *AuthService) VerifyJWTRSA(token string) (bool, *jwt.Token, error) {
	curDir, _ := os.Getwd()
	replaced := strings.Replace(curDir, "/app/middleware", "", 1)
	replaced = strings.Replace(replaced, "\\app\\middleware", "", 1)

	publicKeyPath := config.GetEnv("PUBLIC_KEY_PATH")

	publicKey, err := ioutil.ReadFile(replaced + "/" + publicKeyPath)
	if err != nil {
		return false, nil, errors.New("error reading file public key.pem")
	}
	var parsedToken *jwt.Token

	// parse token
	state, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		// ensure signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unknown signing method")
		}

		parsedToken = token

		// verify
		key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
		if err != nil {
			return nil, err
		}

		return key, nil
	})

	if err != nil {
		return false, &jwt.Token{}, err
	}

	if !state.Valid {
		return false, &jwt.Token{}, errors.New("invalid jwt token")
	}
	return true, parsedToken, nil
}

// UserHasRoles Check if user has roles
func (s *AuthService) UserHasRoles(c *gin.Context, roles ...string) (hasPermission bool) {

	userData, err := helper.DecodedJWTToken(c)
	if err != nil {
		return false
	}
	userRoles := userData["user_roles"]

	for _, role := range roles {
		if Contain(userRoles, role) {
			return true
		}
	}
	return false
}

// UserHasPermissions function to check if this userPermissions contain some permission
func (s *AuthService) UserHasPermissions(c *gin.Context, Permissions ...string) (hasPermission bool) {
	userData, err := helper.DecodedJWTToken(c)
	if err != nil {
		return false
	}
	userPermissions := userData["user_permissions"]

	for _, Permission := range Permissions {
		if Contain(userPermissions, Permission) {
			return true
		}
	}
	return false
}

// Contain function to check if someData contain in someArray
func Contain(someArray interface{}, someData interface{}) bool {
	v := someArray.([]interface{})
	for _, data := range v {
		if data == someData {
			return true
		}
	}
	return false
}
