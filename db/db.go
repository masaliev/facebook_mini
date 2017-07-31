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
	dataStorage.db = db
	return dataStorage
}