package middleware

import (
	"fmt"
	"os"

	"github.com/amteja/ryuk/ecodes"
	"github.com/amteja/ryuk/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	//get token from header
	token := c.GetHeader("Authorization")

	if token == "" {
		c.AbortWithStatusJSON(401, models.ServerResponse{
			Message: "Unauthorized",
			Code:    ecodes.UNAUTHORIZED,
			Data:    nil,
		})
	}

	//validate token

	ok, userId := validateToken(token)

	if !ok {
		c.AbortWithStatusJSON(401, models.ServerResponse{
			Message: "Unauthorized",
			Code:    ecodes.UNAUTHORIZED,
			Data:    nil,
		})
	}

	//set userId in context
	c.Set("userId", userId)

	c.Next()
}

func validateToken(tokenString string) (bool, uint) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return false, 0
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := claims["sub"]

		return true, uint(userId.(float64))
	} else {
		return false, 0
	}

}
