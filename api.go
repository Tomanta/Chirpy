package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/tomanta/chirpy/internal/database"
	"net/http"
	"strings"
	"time"
)

func (cfg *apiConfig) handlerCreateChirp(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	user, err := cfg.dbQueries.GetUser(context.Background(), params.UserID)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "Invalid user_id", err)
	}

	if params.Body == "" {
		respondWithError(writer, http.StatusBadRequest, "Request does not contain body parameter", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(writer, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	newChirp := database.CreateChirpParams{
		Body:   cleanBody(params.Body),
		UserID: user.ID,
	}

	newChirpResponse, err := cfg.dbQueries.CreateChirp(context.Background(), newChirp)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "Could not create chirp", err)
	}

	type createdChirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	// Chirp is under max length
	respondWithJSON(writer, http.StatusCreated, createdChirp{
		ID:        newChirpResponse.ID,
		CreatedAt: newChirpResponse.CreatedAt,
		UpdatedAt: newChirpResponse.UpdatedAt,
		Body:      newChirpResponse.Body,
		UserID:    newChirpResponse.UserID,
	})
}

func (cfg *apiConfig) handlerCreateUser(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	type User struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	// Error: Unable to decode
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	returnUser, err := cfg.dbQueries.CreateUser(context.Background(), params.Email)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	payload := User{
		Id:        returnUser.ID,
		CreatedAt: returnUser.CreatedAt,
		UpdatedAt: returnUser.UpdatedAt,
		Email:     returnUser.Email,
	}

	respondWithJSON(writer, http.StatusCreated, payload)
}

func cleanBody(to_clean string) string {

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(to_clean, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
