package database

// Database interface

type Database interface {
	StoreURL(shortURL, originalURL string) error
	GetURL(shortURL string) (string, error)
}
