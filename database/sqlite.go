package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	db *sql.DB
}

func NewSQLiteDB() *SQLiteDB {
	db, err := sql.Open("sqlite3", "./urls.db")
	if err != nil {
		panic(err)
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS urls (
		short_url TEXT PRIMARY KEY,
		original_url TEXT NOT NULL
	)`) // Create table if not exists
	return &SQLiteDB{db}
}

func (s *SQLiteDB) StoreURL(shortURL, originalURL string) error {
	_, err := s.db.Exec("INSERT INTO urls (short_url, original_url) VALUES (?, ?)", shortURL, originalURL)
	return err
}

func (s *SQLiteDB) GetURL(shortURL string) (string, error) {
	var originalURL string
	err := s.db.QueryRow("SELECT original_url FROM urls WHERE short_url = ?", shortURL).Scan(&originalURL)
	return originalURL, err
}
