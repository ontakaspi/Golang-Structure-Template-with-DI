package main

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
	"time"
)

// Create a custom visitor struct which holds the rate limiter for each
// visitor and the last time that the visitor was seen.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Change the the map to hold values of the type visitor.
var visitors = make(map[string]*visitor)
var mu sync.Mutex

func getVisitor(indexKey string, decodedToken map[string]interface{}) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[indexKey]
	if !exists {
		isBackendService := false
		if decodedToken != nil {
			userRoles := decodedToken["user_roles"].([]interface{})

			//check is the array contains admin

			for _, role := range userRoles {
				if role == "backend_services" {
					isBackendService = true
					break
				}
			}
		}
		var limiter *rate.Limiter
		if isBackendService {
			limiter = rate.NewLimiter(100, 100)
		} else {
			limiter = rate.NewLimiter(5, 30)
		}
		// Include the current time when creating a new visitor.
		visitors[indexKey] = &visitor{limiter, time.Now()}
		return limiter
	}

	// Update the last seen time for the visitor.
	v.lastSeen = time.Now()
	return v.limiter
}

func cleanupVisitors() {
	mu.Lock()
	for indexKey, v := range visitors {
		if time.Since(v.lastSeen) > 3*time.Minute {
			delete(visitors, indexKey)
		}
	}
	mu.Unlock()
}
func decodedJWTToken(authHeader string) (map[string]interface{}, error) {
	const BearerSchema = "Bearer "
	tokenString := authHeader[len(BearerSchema):]
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type errorResp struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}
		var limiter *rate.Limiter
		jwtKey := r.Header.Get("Authorization")
		if jwtKey == "" {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
				return
			}
			limiter = getVisitor(ip, nil)
		} else {
			//convert to gin context
			//c :=
			decodedToken, err := decodedJWTToken(jwtKey)
			if err != nil {
				limiter = getVisitor(jwtKey, nil)
			} else {
				limiter = getVisitor(jwtKey, decodedToken)
			}
		}

		if limiter.Allow() == false {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)

			json.NewEncoder(w).Encode(errorResp{
				Status:  "Too Many Requests",
				Message: "You have made too many requests in a given amount of time. Please try again later.",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
