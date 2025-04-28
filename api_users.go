package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/tomanta/chirpy/internal/auth"
	"github.com/tomanta/chirpy/internal/database"
	"net/http"
	"time"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type UserParameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerCreateUser(writer http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)
	params := UserParameters{}
	err := decoder.Decode(&params)

	// Error: Unable to decode
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	pw_hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Could not hash password", err)
		return
	}

	user_params := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: pw_hash,
	}

	returnUser, err := cfg.dbQueries.CreateUser(context.Background(), user_params)
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

func (cfg *apiConfig) handlerUpdateUser(writer http.ResponseWriter, request *http.Request) {

	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "Could not find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "Could not validate JWT", err)
		return
	}

	decoder := json.NewDecoder(request.Body)
	params := UserParameters{}
	err = decoder.Decode(&params)

	// Error: Unable to decode
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	pw_hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Could not hash password", err)
		return
	}

	user_params := database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: pw_hash,
		ID:             userID,
	}

	updated_user, err := cfg.dbQueries.UpdateUser(context.Background(), user_params)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "Could not hash password", err)
		return
	}

	type Payload struct {
		Updated_at time.Time `json:"updated_at"`
		User_id    uuid.UUID `json:"id"`
		Email      string    `json:"email"`
	}

	respondWithJSON(writer, http.StatusOK, Payload{
		Updated_at: updated_user.UpdatedAt,
		User_id:    updated_user.ID,
		Email:      updated_user.Email,
	})

}
