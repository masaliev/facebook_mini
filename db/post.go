package db


import (
	"github.com/russross/meddler"
	"strconv"
	"time"
)

type Post struct {
	ID int `json:"id" meddler:"id,pk"`
	Text string `json:"text" meddler:"text"`
	UserId int `json:"user_id" meddler:"user_id"`
	CommentsCount int `json:"comments_count" meddler:"comments_count"`
	LikeCount int `json:"like_count" meddler:"like_count"`
	CreateDate time.Time `json:"create_date" meddler:"create_date,localtime"`
}

type PostVerbose struct {
	ID int `json:"id" meddler:"id,pk"`
	Text string `json:"text" meddler:"text"`
	UserId int `json:"user_id" meddler:"user_id"`
	UserName string `json:"user_name" meddler:"user_name"`
	UserPicture string `json:"user_picture" meddler:"user_picture"`
	CommentsCount int `json:"comments_count" meddler:"comments_count"`
	LikeCount int `json:"like_count" meddler:"like_count"`
	CreateDate time.Time `json:"create_date" meddler:"create_date,localtime"`
	LikeId int `json:"like_id" meddler:"like_id"`
}

type PostSortType int

const (
	SortByCreateDate PostSortType = 0
	SortByCommentsCount PostSortType = 1
	SortByLikeCount PostSortType = 2
	PostsPerPage int = 10
)


func (s *DataStorage) SavePost(post *Post) error  {
	if post.ID == 0 {
		return meddler.Insert(s.db, "posts", post)
	}else {
		return meddler.Update(s.db, "posts", post)
	}
}

func (s *DataStorage) GetPosts(sort PostSortType, page int) (error, []*PostVerbose) {
	query := "SELECT p.*, u.full_name AS user_name, u.picture AS user_picture, " +
		"IFNULL((SELECT id FROM likes WHERE post_id = p.id LIMIT 1), 0) AS like_id " +
		"FROM posts AS p LEFT JOIN users AS u ON p.user_id = u.id ";
	switch sort {
	case SortByCreateDate:
		query += " ORDER BY p.create_date "
		break
	case SortByCommentsCount:
		query += " ORDER BY p.comments_count DESC "
		break
	case SortByLikeCount:
		query += " ORDER BY p.like_count DESC"
		break
	}

	query += " LIMIT " + strconv.Itoa(PostsPerPage)
	query += " OFFSET " + strconv.Itoa( page - 1 * PostsPerPage)

	posts := make([]*PostVerbose, 0)
	err := meddler.QueryAll(s.db, &posts, query)
	if err != nil {
		return err, nil
	}

	return nil, posts
}

func (s *DataStorage) DeletePost(post *Post) error {
	_, err := s.db.Exec("DELETE FROM posts WHERE id = ?", post.ID)
	return err
}

func (s *DataStorage) GetPostById(id int) (error, *Post) {
	post := &Post{}
	err := meddler.QueryRow(s.db, post, "SELECT * FROM posts WHERE id = ?", id)
	if err != nil{
		return err, nil
	}
	return nil, post
}