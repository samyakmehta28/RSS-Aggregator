package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/samyakmehta28/RSS-Aggregator/internal/database"
)

func (apiCfg *apiConfig) createFeedHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing request body")
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	respondWithJSON(w, http.StatusCreated, dataBaseFeedToFeed(feed))
}

func (apiCfg *apiConfig) getFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting feeds")
		return
	}

	feedResponses := make([]Feed, len(feeds))
	for i, feed := range feeds {
		feedResponses[i] = dataBaseFeedToFeed(feed)
	}
	respondWithJSON(w, http.StatusOK, feedResponses)

}