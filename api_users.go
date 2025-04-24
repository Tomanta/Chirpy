package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

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
