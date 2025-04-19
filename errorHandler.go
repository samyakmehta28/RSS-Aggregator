package main

import "net/http"

func pathNotFoundError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusNotFound, "Resource not found")
}