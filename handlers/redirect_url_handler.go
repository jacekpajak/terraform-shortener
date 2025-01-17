package handlers

import (
	"net/http"
	"url-shortener/database"
)

func RedirectURLHandler(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := r.URL.Path[1:] // Get path without '/'
		if shortURL == "" {
			http.Error(w, "Short URL required", http.StatusBadRequest)
			return
		}

		// Retrieve the original URL from the database
		originalURL, err := db.GetURL(shortURL)
		if err != nil {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		// Redirect to the original URL
		http.Redirect(w, r, originalURL, http.StatusFound)
	}
}
