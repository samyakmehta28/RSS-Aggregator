package main

import (
	"fmt"
	"net/http"

	"github.com/samyakmehta28/RSS-Aggregator/internal/auth"
	"github.com/samyakmehta28/RSS-Aggregator/internal/database"
)

type authHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		APIKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Invalid API key %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), APIKey)
		
		if err != nil {
			fmt.Printf("API Key: %s", APIKey)
			respondWithError(w, http.StatusInternalServerError, "Error getting user")
			return
		}

		handler(w,r,user)
	}
}