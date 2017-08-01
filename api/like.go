package api

import (
	"github.com/labstack/echo"
	"strconv"
	"net/http"
)

func (a *Api) Like (c echo.Context) error {
	userId := GetUserIDFromToken(c)
	postId,_ := strconv.Atoi(c.FormValue("post_id"))

	if postId == 0{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Bad Request"}
	}

	isLiked := a.dataStorage.IsLiked(userId, postId)

	if isLiked{
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Allready liked"}
	}

	err, post := a.dataStorage.GetPostById(postId)
	if err != nil || (post != nil && post.ID == 0){
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Post not found"}
	}

	err, like := a.dataStorage.Like(userId, post)
	if err != nil{
		return err
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

	err, post := a.dataStorage.GetPostById(postId)
	if err != nil || (post != nil && post.ID == 0){
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Post not found"}
	}

	err = a.dataStorage.Unlike(currentUserId, post)
	if err != nil{
		return err
	}

	return c.NoContent(http.StatusOK)
}

