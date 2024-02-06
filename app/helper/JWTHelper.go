package helper

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// DecodedToken Decode a JWT token without validation for getting payload data (claim)
func DecodedJWTToken(c *gin.Context) (map[string]interface{}, error) {

	const BearerSchema = "Bearer "
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header is required")
	}
	tokenString := authHeader[len(BearerSchema):]
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil

}
