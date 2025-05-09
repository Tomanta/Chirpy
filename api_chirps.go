package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/tomanta/chirpy/internal/auth"
	"github.com/tomanta/chirpy/internal/database"
	"net/http"
	"strings"
	"time"
	"sort"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerGetChirps(writer http.ResponseWriter, request *http.Request) {

	author := request.URL.Query().Get("author_id")
	authorID := uuid.Nil
	if author != "" {
		var err error
		authorID, err = uuid.Parse(author)
		if err != nil {
			respondWithError(writer, http.StatusBadRequest, "Invalid author id", err)
			return
		}
	}

	sortOrder := request.URL.Query().Get("sort")
	if sortOrder != "" {
		if sortOrder != "asc" && sortOrder != "desc" {
			respondWithError(writer, http.StatusBadRequest, "Invalid sort order", nil)
			return
		}
	}

	dbChirps, err := cfg.dbQueries.GetChirps(context.Background())
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "Could not retrieve chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, c := range dbChirps {
		if author == "" || authorID == c.UserID {
			chirps = append(chirps, Chirp{
				ID:        c.ID,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
				Body:      c.Body,
				UserID:    c.UserID,
			})
		}
	}

	if sortOrder == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	respondWithJSON(writer, http.StatusOK, chirps)

}

func (cfg *apiConfig) handlerGetChirpByID(writer http.ResponseWriter, request *http.Request) {

	chirpID, err := uuid.Parse(request.PathValue("chirpID"))
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbResponse, err := cfg.dbQueries.GetChirpByID(context.Background(), chirpID)
	if err != nil {
		respondWithError(writer, http.StatusNotFound, "Could not retrieve chirp", err)
		return
	}

	respondWithJSON(writer, http.StatusOK, Chirp{
		ID:        dbResponse.ID,
		CreatedAt: dbResponse.CreatedAt,
		UpdatedAt: dbResponse.UpdatedAt,
		Body:      dbResponse.Body,
		UserID:    dbResponse.UserID,
	})

}

func (cfg *apiConfig) handlerCreateChirp(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

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
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Could not decode parameters", err)
		return
	}

	// user, err := cfg.dbQueries.GetUser(context.Background(), params.UserID)
	// if err != nil {
	// 	respondWithError(writer, http.StatusBadRequest, "Invalid user_id", err)
	// 	return
	// }

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
		UserID: userID,
	}

	newChirpResponse, err := cfg.dbQueries.CreateChirp(context.Background(), newChirp)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "Could not create chirp", err)
		return
	}

	// Chirp is under max length
	respondWithJSON(writer, http.StatusCreated, Chirp{
		ID:        newChirpResponse.ID,
		CreatedAt: newChirpResponse.CreatedAt,
		UpdatedAt: newChirpResponse.UpdatedAt,
		Body:      newChirpResponse.Body,
		UserID:    newChirpResponse.UserID,
	})
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

func (cfg *apiConfig) handlerDeleteChirpByID(writer http.ResponseWriter, request *http.Request) {
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

	chirpID, err := uuid.Parse(request.PathValue("chirpID"))
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbResponse, err := cfg.dbQueries.GetChirpByID(context.Background(), chirpID)
	if err != nil {
		respondWithError(writer, http.StatusNotFound, "Could not retrieve chirp", err)
		return
	}

	if dbResponse.UserID != userID {
		respondWithError(writer, http.StatusForbidden, "", nil)
		return
	}

	err = cfg.dbQueries.DeleteChirpByID(context.Background(), dbResponse.ID)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Could not delete chirp", err)
		return
	}

	writer.WriteHeader(http.StatusNoContent)

}
