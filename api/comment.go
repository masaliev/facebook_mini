package api

import (
	"github.com/labstack/echo"
	"strconv"
	"net/http"
)

func (a *Api) GetComments(c echo.Context) error {
	userId := GetUserIDFromToken(c)
	postId, _ := strconv.Atoi(c.Param("postId"))

	err, comments := a.dataStorage.GetComments(userId, postId)
	if err != nil{
		return err
	}
	return c.JSON(http.StatusOK, comments)
}

func (a *Api)AddComment(c echo.Context) error {
	userId := GetUserIDFromToken(c)
	postId, _ := strconv.Atoi(c.Param("postId"))
	text := c.FormValue("text")

	err, post := a.dataStorage.GetPostById(postId)
	if err != nil || (post != nil && post.ID == 0){
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Post not found"}
	}

	err, comment := a.dataStorage.AddComment(userId, post, text)
	if err != nil{
		return err
	}

	err, commentDetail := a.dataStorage.GetCommentDetailById(comment.ID)

	if err != nil{
		 return err
	}
	return c.JSON(http.StatusCreated, commentDetail)
}

func (a *Api) UpdateComment(c echo.Context) error  {
	userId := GetUserIDFromToken(c)
	commentId,_ := strconv.Atoi(c.Param("id"))
	text := c.FormValue("text")

	err, comment := a.dataStorage.GetCommentById(commentId)
	if err != nil || (comment != nil && comment.ID == 0){
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Comment not found"}
	}

	if userId != comment.UserId{
		return &echo.HTTPError{Code: http.StatusForbidden, Message: "Forbidden"}
	}

	comment.Text = text

	err = a.dataStorage.SaveComment(comment)
	if err != nil{
		return err
	}

	err, commentDetail := a.dataStorage.GetCommentDetailById(comment.ID)
	if err != nil{
		return err
	}
	return c.JSON(http.StatusOK, commentDetail)
}

func (a *Api) DeleteComment(c echo.Context) error {
	userId := GetUserIDFromToken(c)
	commentId, _ := strconv.Atoi(c.Param("id"))

	err, comment := a.dataStorage.GetCommentById(commentId)
	if err != nil || (comment != nil && comment.ID == 0){
		return err
	}

	if comment.UserId != userId{
		return &echo.HTTPError{Code: http.StatusForbidden, Message: "Forbiden"}
	}

	if err := a.dataStorage.DeleteComment(comment); err != nil{
		return err
	}
	return c.NoContent(http.StatusOK)
}
