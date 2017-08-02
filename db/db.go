package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type (
	DataStorage struct {
		db *sql.DB
	}
)

func NewDataStorage(dbPath string) *DataStorage {
	dataStorage := new(DataStorage)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil{
		panic(err)
	}

	//Init tables
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, full_name VARCHAR(255), password VARCHAR(500), picture VARCHAR(255), phone varchar(20));")
	if err != nil{
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY, text TEXT, user_id INTEGER, comments_count INTEGER, like_count INTEGER, create_date DATETIME DEFAULT CURRENT_TIMESTAMP);")
	if err != nil{
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS likes (id INTEGER PRIMARY KEY, user_id INTEGER, post_id INTEGER, create_date DATETIME DEFAULT CURRENT_TIMESTAMP);")
	if err != nil{
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS comments (id INTEGER PRIMARY KEY, text TEXT, post_id INTEGER, user_id INTEGER, create_date DATETIME DEFAULT CURRENT_TIMESTAMP);")
	if err != nil{
		panic(err)
	}

	dataStorage.db = db
	return dataStorage
}