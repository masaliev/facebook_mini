package api

import (
	"github.com/labstack/echo"
	"strconv"
	"github.com/masaliev/facebook_mini/db"
	"net/http"
)

func (a *Api) Like (c echo.Context) error {
	userId := GetUserIDFromToken(c)

	l := &db.Like{}
	if err := c.Bind(l); err != nil{
		return err
	}

	if l.UserId != userId{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Bad Request"}
	}

	err, isLiked := a.dataStorage.IsLiked(l.UserId, l.PostId)
	if err != nil{
		return err
	}

	if isLiked{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Allready like"}
	}
	err, like := a.dataStorage.Like(l.UserId, l.PostId)
	if err != nil{
		return err
	}

	return c.JSON(http.StatusOK, like)
}

func (a *Api) UnLike (c echo.Context) error {
	currentUserId := GetUserIDFromToken(c)
	userId,_ := strconv.Atoi(c.Param("user_id"))
	postId,_ := strconv.Atoi(c.Param("post_id"))

	if userId == 0 || postId == 0 || currentUserId != userId{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Bad Request"}
	}

	err, isLiked := a.dataStorage.IsLiked(userId, postId)
	if err != nil{
		return err
	}

	if !isLiked{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Like not found"}
	}
	err = a.dataStorage.Unlike(userId, postId)
	if err != nil{
		return err
	}

	return c.NoContent(http.StatusOK)
}

