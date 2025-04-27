package main

import (
	"context"
	"encoding/json"
	"github.com/tomanta/chirpy/internal/auth"
	"github.com/tomanta/chirpy/internal/database"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerUserLogin(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
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

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't create JWT token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	_, err = cfg.dbQueries.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}

	payload := response{
		User: User{
			Id:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	}

	respondWithJSON(writer, http.StatusOK, payload)

}

func (cfg *apiConfig) handlerRefresh(writer http.ResponseWriter, request *http.Request) {
	refreshToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "Couldn't find token", err)
	}

	user, err := cfg.dbQueries.GetUserFromRefreshToken(context.Background(), refreshToken)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "Couldn't get user from refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}
	respondWithJSON(writer, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(writer http.ResponseWriter, request *http.Request) {
	refreshToken, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "Couldn't find token", err)
	}

	err = cfg.dbQueries.RevokeRefreshToken(context.Background(), refreshToken)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Could not revoke session", err)
	}

	writer.WriteHeader(http.StatusNoContent)

}
