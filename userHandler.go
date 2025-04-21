package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/samyakmehta28/RSS-Aggregator/internal/database"
)

func (apiCfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing request body")
		return
	}

	user, err:= apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:  uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: params.Name,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	respondWithJSON(w, http.StatusCreated, dataBaseUserToUser(user))
}


func (apiCfg *apiConfig) getUserByAPIKey(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, dataBaseUserToUser(user))
}


// type GetPostForUserParams struct {
// 	UserID uuid.UUID
// 	Limit  int32
// }


func (apiCfg *apiConfig) getPostForUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	posts, err := apiCfg.DB.GetPostForUser(r.Context(), database.GetPostForUserParams{
		UserID: user.ID,
		Limit:  10,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching posts for user")
		return
	}

	postsResponse := make([]Post, len(posts))
	for i, post := range posts {
		postsResponse[i] = dataBasePostToPost(post)
	respondWithJSON(w, http.StatusOK, dataBaseUserToUser(user))
	}
}
