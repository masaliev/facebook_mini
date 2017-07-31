package db


import (
	"github.com/russross/meddler"
)

type Like struct {
	ID int `json:"id" meddler:"id,pk"`
	PostId int `json:"post_id" meddler:"post_id"`
	UserId int `json:"user_id" meddler:"user_id"`
	CreateDate int `json:"create_date" meddler:"create_date"`
}

func (s *DataStorage) IsLiked(userId, postId int) (error, bool) {
	var l Like
	err := meddler.QueryRow(s.db, l, "SELECT * FROM likes WHERE user_id = ? and post_id = ?", userId, postId)
	if err != nil{
		return err, false
	}

	if l.ID != 0{
		return nil, true
	}else{
		return nil, false
	}
}

func (s *DataStorage) Like(userId, postId int) (error, *Like ){
	like := &Like{
		PostId: postId,
		UserId: userId,
	}
	err := meddler.Insert(s.db, "likes", like)
	if err != nil{
		return err, nil
	}
	return nil, like
}

func (s *DataStorage) Unlike(userId, postId int) error {
	_, err := s.db.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userId, postId)
	return err
}