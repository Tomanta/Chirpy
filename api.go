package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func handlerValidateChirp(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnClean struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	// Error: Unable to decode
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Body == "" {
		respondWithError(writer, http.StatusBadRequest, "Request does not contain body parameter", err)
		return
	}

	const maxChirpLength = 140
	// Error: Chirp is too long
	if len(params.Body) > maxChirpLength {
		respondWithError(writer, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	// Chirp is under max length
	respondWithJSON(writer, http.StatusOK, returnClean{
		CleanedBody: cleanBody(params.Body),
	})

}

func cleanBody(to_clean string) string {
	words := strings.Split(to_clean, " ")
	for i, word := range words {
		lower := strings.ToLower(word)
		if lower == "kerfuffle" || lower == "sharbert" || lower == "fornax" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
