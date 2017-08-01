package db


import (
	"github.com/russross/meddler"
	"time"
)

type Like struct {
	ID int `json:"id" meddler:"id,pk"`
	PostId int `json:"post_id" meddler:"post_id"`
	UserId int `json:"user_id" meddler:"user_id"`
	CreateDate time.Time `json:"create_date" meddler:"create_date,localtime"`
}

func (s *DataStorage) IsLiked(userId, postId int) bool {
	l := &Like{}
	err := meddler.QueryRow(s.db, l, "SELECT * FROM likes WHERE user_id = ? and post_id = ?", userId, postId)
	if err != nil{
		return false
	}

	if l.ID != 0{
		return true
	}else{
		return false
	}
}

func (s *DataStorage) Like(userId int, post *Post) (error, *Like ){
	like := &Like{
		PostId: post.ID,
		UserId: userId,
		CreateDate: time.Now(),
	}
	err := meddler.Insert(s.db, "likes", like)
	if err != nil{
		return err, nil
	}

	post.LikeCount++
	s.SavePost(post)

	return nil, like
}

func (s *DataStorage) Unlike(userId int, post *Post) error {
	_, err := s.db.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userId, post.ID)
	if err != nil{
		return err
	}

	post.LikeCount--
	s.SavePost(post)

	return err
}
