package db

import "github.com/russross/meddler"

type User struct {
	ID       int `json:"id" meddler:"id,pk"`
	FullName string `json:"full_name" meddler:"full_name"`
	Picture  string `json:"picture" meddler:"picture"`
	Phone    string `json:"phone" meddler:"phone"`
	Password string `json:"-" meddler:"password"`
	Token    string `json:"token" meddler:"-"`
}


func (s *DataStorage) SaveUser(user *User) error {
	return meddler.Insert(s.db, "users", user)
}

func (s *DataStorage) GetByPhone(phone string) (error, *User) {
	user := new(User)
	err := meddler.QueryRow(s.db, user, "select * from users where phone = ?", phone)
	return err, user
}

func (s *DataStorage) UpdateUser(user *User) error {
	return meddler.Update(s.db, "users", user)
}
