package api


import (
	"github.com/labstack/echo"
	"strconv"
	"net/http"
	"github.com/masaliev/facebook_mini/db"
	"time"
)

type (
	Post struct {
		Text string `json:"text"`
	}
)

func (a *Api) GetPosts(c echo.Context) error {
	//userId := GetUserIDFromToken(c)
	page, _ := strconv.Atoi(c.QueryParam("page"))
	sort := c.QueryParam("sort")

	if page == 0{
		page = 1
	}

	var sortType db.PostSortType
	switch sort {
	case "comments":
		sortType = db.SortByCommentsCount
		break
	case "likes":
		sortType = db.SortByLikeCount
		break
	default:
		sortType = db.SortByCreateDate
	}
	err, posts := a.dataStorage.GetPosts(sortType, page)
	if err != nil{
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	return c.JSON(http.StatusOK, posts)
}

func (a *Api) GetPost(c echo.Context) error {
	postId,_ := strconv.Atoi(c.Param("id"))
	err, post := a.dataStorage.GetPostDetailById(postId)
	if err != nil{
		return err
	}
	return c.JSON(http.StatusOK, post)
}

func (a *Api) CreatePost(c echo.Context) error {
	userId := GetUserIDFromToken(c)

	p := &Post{}
	if err := c.Bind(p); err != nil{
		return err
	}

	if p.Text == ""{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Text is empty"}
	}

	post := &db.Post{
		Text: p.Text,
		UserId: userId,
		CreateDate: time.Now(),
	}

	if err := a.dataStorage.SavePost(post); err != nil{
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Try again"}
	}

	return c.JSON(http.StatusCreated, post)
}

func (a *Api) UpdatePost(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Bad id"}
	}

	userId := GetUserIDFromToken(c)

	p := &Post{}
	if err := c.Bind(p); err != nil{
		return err
	}

	if p.Text == ""{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Text is empty"}
	}

	err, post := a.dataStorage.GetPostById(id)
	if err != nil{
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Try again"}
	}
	if post == nil || post.ID == 0{
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "Not found"}
	}

	if post.UserId != userId{
		return &echo.HTTPError{Code: http.StatusForbidden, Message: "Forbiden"}
	}

	post.Text = p.Text

	if err := a.dataStorage.SavePost(post); err != nil{
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Try again"}
	}

	return c.JSON(http.StatusOK, post)
}

func (a *Api) DeletePost(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Bad id"}
	}

	userId := GetUserIDFromToken(c)

	err, post := a.dataStorage.GetPostById(id)
	if err != nil{
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Try again"}
	}
	if post == nil || post.ID == 0{
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "Not found"}
	}

	if post.UserId != userId{
		return &echo.HTTPError{Code: http.StatusForbidden, Message: "Forbiden"}
	}

	if err := a.dataStorage.DeletePost(post); err != nil{
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "Try again"}
	}
	return c.NoContent(http.StatusOK)
}
