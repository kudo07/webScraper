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
