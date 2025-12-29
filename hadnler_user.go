package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kudo07/webScraper/internal/database"
)

func (apiCnfg *apiConfig) hadnlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Name string `json:"name"`
	}
	p := params{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payloads")
		return
	}
	user, err := apiCnfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      p.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could no create user: %s", err.Error()))
		return
	}
	respondWithJSON(w, http.StatusCreated, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	// authenticated only user
	// apiKey, err := auth.GetApiKey(r.Header)
	// if err != nil {
	// 	respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
	// 	return
	// }
	// user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
	// if err != nil {
	// 	respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
	// 	return
	// }
	// respondWithJSON(w, 200, databaseUserToUser(user))

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
