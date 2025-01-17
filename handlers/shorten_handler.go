package handlers

import (
	"encoding/json"
	"net/http"
	"url-shortener/database"
	"url-shortener/utils"
)

func StoreURLHandler(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		originalURL := body["url"]
		if originalURL == "" {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		shortURL := utils.GenerateShortURL()
		err := db.StoreURL(shortURL, originalURL)
		if err != nil {
			http.Error(w, "Error storing URL", http.StatusInternalServerError)
			return
		}

		resp := map[string]string{"short_url": shortURL, "original_url": originalURL}
		json.NewEncoder(w).Encode(resp)
	}
}
