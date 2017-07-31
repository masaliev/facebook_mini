package api

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/labstack/echo"
	"strconv"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"github.com/masaliev/facebook_mini/db"
)

type(
	User struct {
		FullName string `json:"full_name,omitempty"`
		Phone string `json:"phone"`
		Password string `json:"password"`
	}
)

func (a *Api) SignUp(c echo.Context) error {
	u := &User{}
	if err := c.Bind(u); err != nil{
		return err
	}

	var message string
	if u.FullName == ""{
		message = "Full name is empty"
	} else if u.Phone == ""{
		message = "Phone is empty"
	}else if u.Password == ""{
		message = "Password is empty"
	}

	if message != ""{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: message}
	}

	password, err := generatePasswordHash(u.Password)
	if err != nil{
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Try again"}
	}

	user := &db.User{
		FullName: u.FullName,
		Phone: u.Phone,
		Password: string(password),
	}

	if err := a.dataStorage.SaveUser(user); err != nil{
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Try again"}
	}

	token, err := createJWTToken(user)
	if err != nil{
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Try again"}
	}
	user.Token = token

	return c.JSON(http.StatusCreated, user)
}

func (a *Api) Login(c echo.Context) error {
	u := &User{}
	if err := c.Bind(u); err != nil{
		return err
	}

	if u.Phone == "" || u.Password == ""{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message:"Invalid phone or password"}
	}

	err, user := a.dataStorage.GetByPhone(u.Phone)
	if err != nil || user.ID == 0{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message:"Invalid phone or password"}
	}

	if err := compareHashAndPassword(user.Password, u.Password); err != nil{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message:"Invalid phone or password"}
	}

	token, err := createJWTToken(user)
	if err != nil{
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Try again"}
	}
	user.Token = token

	return c.JSON(http.StatusOK, user)
}

func (a *Api) LoginTest(c echo.Context) error {
	userId := GetUserIDFromToken(c)

	return c.String(http.StatusOK, strconv.Itoa(userId))
}

func generatePasswordHash(password string) ([]byte, error){
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func compareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func createJWTToken(user *db.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	return token.SignedString([]byte(Key))
}


func GetUserIDFromToken(c echo.Context) int {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	id := claims["id"].(float64)

	return int(id)
}