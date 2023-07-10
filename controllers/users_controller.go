package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/amteja/ryuk/ecodes"
	"github.com/amteja/ryuk/intializers"
	"github.com/amteja/ryuk/models"
	"github.com/amteja/ryuk/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	//get email and pass
	var body struct {
		Email    string
		Password string
	}

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Message: "Bad Body Format",
			Code:    ecodes.EMAIL_AND_PASSWORD_FORMAT_INVALID,
			Data:    nil,
		})
		return
	}

	//validate
	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Message: "Email and Password are required",
			Code:    ecodes.EMAIL_AND_PASSWORD_REQUIRED,
			Data:    nil,
		})
		return
	}

	if !validateEmail(body.Email) {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Message: "Invalid Email Format",
			Code:    ecodes.INVALID_EMAIL_FORMAT,
			Data:    nil,
		})
		return
	}

	if code := validatePassword(body.Password); code != ecodes.NO_ERROR {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Message: "Invalid Password Format",
			Code:    code,
			Data:    nil,
		})
		return
	}

	//hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Message: "Something went wrong",
			Code:    ecodes.INTERNAL_SERVER_ERROR,
			Data:    nil,
		})
		return
	}

	//create user
	user := models.User{
		Email:    body.Email,
		Password: string(hash),
	}

	result := intializers.DB.Create(&user)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "violates unique") {
			c.JSON(http.StatusBadRequest, models.ServerResponse{
				Message: "User already exists",
				Code:    ecodes.USER_ALREADY_EXISTS,
				Data:    nil,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Message: "Something went wrong",
			Code:    ecodes.INTERNAL_SERVER_ERROR,
			Data:    nil,
		})
		return
	}

	//send data with user password removed
	c.JSON(http.StatusOK, models.ServerResponse{
		Message: "User created successfully",
		Code:    ecodes.NO_ERROR,
		Data:    user,
	})

	//return user
}

func Login(c *gin.Context) {

	var body struct {
		Email    string
		Password string
	}

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Message: "Bad Body Format",
			Code:    ecodes.EMAIL_AND_PASSWORD_FORMAT_INVALID,
			Data:    nil,
		})
		return
	}

	//validate
	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Message: "Email and Password are required",
			Code:    ecodes.EMAIL_AND_PASSWORD_REQUIRED,
			Data:    nil,
		})
		return
	}

	if !validateEmail(body.Email) {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Message: "Invalid Credentials",
			Code:    ecodes.INCORRECT_CREDENTIALS,
			Data:    nil,
		})
		return
	}

	if code := validatePassword(body.Password); code != ecodes.NO_ERROR {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Message: "Invalid Credentials",
			Code:    ecodes.INCORRECT_CREDENTIALS,
			Data:    nil,
		})
		return
	}

	//get user
	var user models.User
	result := intializers.DB.Where("email = ?", body.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Message: "Invalid Credentials",
			Code:    ecodes.INCORRECT_CREDENTIALS,
			Data:    nil,
		})
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Message: "Invalid Credentials",
			Code:    ecodes.INCORRECT_CREDENTIALS,
			Data:    nil,
		})
		return
	}

	//compare password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, models.ServerResponse{
			Message: "Invalid Credentials",
			Code:    ecodes.INCORRECT_CREDENTIALS,
			Data:    nil,
		})
		return
	}

	//generate token
	token, err := utils.GenerateToken(user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ServerResponse{
			Message: "Something went wrong",
			Code:    ecodes.INTERNAL_SERVER_ERROR,
			Data:    nil,
		})
		return
	}

	//send data with user password removed
	c.JSON(http.StatusOK, models.ServerResponse{
		Message: "User logged in successfully",
		Code:    ecodes.NO_ERROR,
		Data:    gin.H{"token": token},
	})

}

func validateEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
	return emailRegex.MatchString(email)
}

func validatePassword(password string) int {

	if len(password) < 8 {
		return ecodes.PASSWORD_LENGTH_TOO_SHORT
	} else if len(password) > 32 {
		return ecodes.PASSWORD_LENGTH_TOO_LONG
	}

	var passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+=-]+$`)
	if !passwordRegex.MatchString(password) {
		return ecodes.INVALID_PASSWORD_FORMAT
	}
	return ecodes.NO_ERROR
}
