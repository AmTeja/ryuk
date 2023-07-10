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
	if !validateToken(token) {
		c.AbortWithStatusJSON(401, models.ServerResponse{
			Message: "Unauthorized",
			Code:    ecodes.UNAUTHORIZED,
			Data:    nil,
		})
	}

	c.Next()
}

func validateToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["sub"])
	} else {
		fmt.Println(err)
	}

	return true

}
