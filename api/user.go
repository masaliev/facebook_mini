package api

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/labstack/echo"
	"strconv"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"github.com/masaliev/facebook_mini/db"
	"os"
	"io"
)

type(
	User struct {
		FullName string `json:"full_name,omitempty" form:"full_name"`
		Phone string `json:"phone" form:"phone"`
		Password string `json:"password" form:"password"`
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
		return err
	}

	user := &db.User{
		FullName: u.FullName,
		Phone: u.Phone,
		Password: string(password),
	}

	if err := a.dataStorage.SaveUser(user); err != nil{
		return err
	}

	token, err := createJWTToken(user)
	if err != nil{
		return err
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
	if err != nil{
		return err
	}
	if user == nil || user.ID == 0{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message:"Invalid phone or password"}
	}

	if err := compareHashAndPassword(user.Password, u.Password); err != nil{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message:"Invalid phone or password"}
	}

	token, err := createJWTToken(user)
	if err != nil{
		return err
	}
	user.Token = token

	return c.JSON(http.StatusOK, user)
}

func (a *Api) LoginTest(c echo.Context) error {
	userId := GetUserIDFromToken(c)

	return c.String(http.StatusOK, strconv.Itoa(userId))
}

func (a *Api) UploadPicture(c echo.Context) error {
	userId := GetUserIDFromToken(c)

	err, user := a.dataStorage.GetUserById(userId)
	if err != nil{
		return err
	}

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	contentType := file.Header["Content-Type"][0]
	if contentType != "image/jpeg" && contentType != "image/png"  && contentType != "image/gif"{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "The file is not an image file"}
	}

	// Destination
	path := "/avatars/" + strconv.FormatInt(time.Now().Unix(),10) + "_" + file.Filename
	dst, err := os.Create("public" + path)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	if user.Picture != ""{
		os.Remove("public" + user.Picture)
	}

	user.Picture = path
	err = a.dataStorage.UpdateUser(user)
	if err != nil{
		return err
	}

	token, err := createJWTToken(user)
	if err == nil{
		user.Token = token
	}

	return c.JSON(http.StatusOK, user)
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
