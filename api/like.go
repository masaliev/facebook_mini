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

	isLiked := a.dataStorage.IsLiked(l.UserId, l.PostId)

	if isLiked{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Allready liked"}
	}
	err, like := a.dataStorage.Like(l.UserId, l.PostId)
	if err != nil{
		return err
	}

	err, post := a.dataStorage.GetPostById(l.PostId)
	if err == nil{
		post.LikeCount++
		a.dataStorage.SavePost(post)
	}

	return c.JSON(http.StatusOK, like)
}

func (a *Api) UnLike (c echo.Context) error {
	currentUserId := GetUserIDFromToken(c)
	postId,_ := strconv.Atoi(c.Param("post_id"))

	if postId == 0{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Bad Request"}
	}

	isLiked := a.dataStorage.IsLiked(currentUserId, postId)
	if !isLiked{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Like not found"}
	}
	err := a.dataStorage.Unlike(currentUserId, postId)
	if err != nil{
		return err
	}

	err, post := a.dataStorage.GetPostById(postId)
	if err == nil{
		post.LikeCount--
		a.dataStorage.SavePost(post)
	}

	return c.NoContent(http.StatusOK)
}

