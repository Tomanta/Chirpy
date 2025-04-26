package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
	"github.com/tomanta/chirpy/internal/auth"
	"github.com/tomanta/chirpy/internal/database"	
)

type User struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type UserParameters struct {
	Email string `json:"email"`
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
		Email: params.Email,
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


func (cfg *apiConfig) handlerUserLogin(writer http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)
	params := UserParameters{}
	err := decoder.Decode(&params)

	// Error: Unable to decode
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "ncorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}


	payload := User{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(writer, http.StatusOK, payload)	

}