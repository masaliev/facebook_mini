package db

import (
	"time"
	"github.com/russross/meddler"
)

type (
	Comment struct {
		ID int `json:"id" meddler:"id,pk"`
		Text string `json:"text" meddler:"text"`
		PostId int `json:"post_id" meddler:"post_id"`
		UserId int `json:"user_id" meddler:"user_id"`
		CreateDate time.Time `json:"create_date" meddler:"create_date,localtime"`
	}
	CommentDetail struct {
		ID int `json:"id" meddler:"id,pk"`
		Text string `json:"text" meddler:"text"`
		UserId int `json:"user_id" meddler:"user_id"`
		UserName string `json:"user_name" meddler:"user_name"`
		UserPicture string `json:"user_picture" meddler:"user_picture"`
		CreateDate time.Time `json:"create_date" meddler:"create_date,localtime"`
	}
)

func (s *DataStorage) AddComment(userId int, post *Post, text string) (error, *Comment) {
	comment := &Comment{
		Text: text,
		PostId: post.ID,
		UserId: userId,
		CreateDate: time.Now(),
	}

	err := meddler.Insert(s.db, "comments", comment)
	if err != nil{
		return err, nil
	}

	post.CommentsCount++
	s.SavePost(post)

	return nil, comment
}

func (s *DataStorage) SaveComment(comment *Comment) error {
	return meddler.Update(s.db, "comments", comment)
}

func (s *DataStorage) GetComments(userId, postId int) (error, []*CommentDetail) {
	query := "SELECT c.*, u.full_name AS user_name, u.picture AS user_picture FROM comments AS c " +
		"LEFT JOIN users AS u ON c.user_id = u.id "
	comments := make([]*CommentDetail, 0)
	err := meddler.QueryAll(s.db, &comments, query)
	if err != nil{
		return err, nil
	}
	return nil, comments
}

func (s *DataStorage) GetCommentById(id int) (error, *Comment) {
	comment := &Comment{}
	err := meddler.QueryRow(s.db, comment, "SELECT * FROM comments WHERE id = ?", id)
	if err != nil{
		return err, nil
	}
	return nil, comment
}

func (s *DataStorage) GetCommentDetailById(id int) (error, *CommentDetail) {
	comment := &CommentDetail{}
	query := "SELECT c.*, u.full_name AS user_name, u.picture AS user_picture FROM comments AS c " +
		"LEFT JOIN users AS u ON c.user_id = u.id WHERE c.id = ?"

	err := meddler.QueryRow(s.db, comment, query, id)
	if err != nil{
		return err, nil
	}
	return nil, comment
}


func (s *DataStorage) DeleteComment(comment *Comment) error {
	_, err := s.db.Exec("DELETE FROM comments WHERE id = ?", comment.ID)
	if err != nil{
		return err
	}

	err, post := s.GetPostById(comment.PostId)
	if err == nil{
		post.CommentsCount--
		s.SavePost(post)
	}
	return nil
}



