package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
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


func (apiCfg *apiConfig) createFeedFollowUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing request body")
		return
	}

	feedFollowUser, err := apiCfg.DB.CreateFeedFollowUser(r.Context(), database.CreateFeedFollowUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating feed follow user %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, dataBaseFeedsUserToFeedsUser(feedFollowUser))
}


func (apiCfg *apiConfig) getFeedsFollowUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	feeds, err := apiCfg.DB.GetFeedFollowUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting feeds %v", err))
		return
	}

	feedsFollowedByUser := make([]FeedsUser, len(feeds))
	for i, feed := range feeds {
		feedsFollowedByUser[i] = dataBaseFeedsUserToFeedsUser(feed)
	}
	respondWithJSON(w, http.StatusOK, feedsFollowedByUser)
}

func (apiCfg *apiConfig) deleteFeedFollowUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	
	feedFollowUserIDString := chi.URLParam(r, "feedFollowUserID")
	feedFollowUserID, err := uuid.Parse(feedFollowUserIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feed follow user ID")
		return
	}


	err = apiCfg.DB.DeleteFeedFollowUser(r.Context(), database.DeleteFeedFollowUserParams{
		UserID: user.ID,
		FeedID: feedFollowUserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting feed follow user %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Feed follow user deleted successfully"})

}