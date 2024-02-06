package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/goioc/di"
	"github.com/golang-jwt/jwt/v4"
	"golang-structure-template-with-di/app/service"
	"golang-structure-template-with-di/config"
	"net/http"
	"strings"
)

// Middleware for Auth JWT
func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		AuthServiceInstance, err := di.GetInstanceSafe("AuthService")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": err.Error()})
			return
		}
		AuthService := AuthServiceInstance.(*service.AuthService)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "No Authorization header found"})
			return
		}
		const BearerSchema = "Bearer "
		//
		if !strings.Contains(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "A Bearer token must be set"})
			return
		}
		tokenString := authHeader[len(BearerSchema):]
		valid, token, err := AuthService.VerifyJWTRSA(tokenString)
		if valid {
			claims := token.Claims.(jwt.MapClaims)
			if claims["iss"] != config.GetEnv("JWT_ISS") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Issuer not valid"})
				return
			}
		} else {
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": err.Error()})
				return
			}
		}

	}
}
