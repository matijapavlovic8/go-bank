package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	. "go-bank-v2/internal/infrastructure/postgres"
	. "go-bank-v2/internal/types"
	"net/http"
	"os"
	"strconv"
)

func withJWTAuth(handlerFunc gin.HandlerFunc, s Store, adminOnly bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("calling JWT auth middleware")

		tokenString := c.GetHeader("x-jwt-token")
		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(c)
			return
		}
		if !token.Valid {
			permissionDenied(c)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		if checkIfAdminToken(claims["ownerId"].(float64), s) {
			handlerFunc(c)
			return
		}

		if adminOnly && !checkIfAdminToken(claims["ownerId"].(float64), s) {
			permissionDenied(c)
			return
		}

		userID, err := getIDFromContext(c)
		if err != nil {
			permissionDenied(c)
			return
		}
		user, err := s.GetUserByID(userID)
		if err != nil || user == nil {
			permissionDenied(c)
			return
		}

		if user.ID != int(claims["ownerId"].(float64)) {
			permissionDenied(c)
			return
		}

		handlerFunc(c)
	}
}

func permissionDenied(c *gin.Context) {
	c.JSON(http.StatusForbidden, Error{Error: "Permission denied"})
	c.Abort()
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func createJWT(user *User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"ownerId":   user.ID,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

// This function extracts the user ID from the Gin context.
func getIDFromContext(c *gin.Context) (int, error) {
	idStr := c.Param("id")
	if idStr == "" {
		idStr = c.Query("ownerId")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid ID")
	}
	return id, nil
}

func checkIfAdminToken(id float64, s Store) bool {

	user, err := s.GetUserByID(int(id))
	if err != nil {
		return false
	}
	if user.Role == AdminRole {
		return true
	}
	return false
}
