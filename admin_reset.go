package main

import (
	"context"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("DEBUG: Resetting...")

	if cfg.platform != "dev" {
		respondWithError(writer, http.StatusForbidden, "", nil)
		return
	}

	type Reset struct {
		HitCount   int `json:"hit_count"`
		UserCount  int `json:"user_count"`
		ChirpCount int `json:"chirp_count"`
	}

	cfg.fileserverHits.Store(0)

	err := cfg.dbQueries.ResetUsers(context.Background())
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't reset users", err)
		return
	}

	err = cfg.dbQueries.ResetChirps(context.Background())
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't reset chirps", err)
		return
	}

	payload := Reset{
		HitCount:   0,
		UserCount:  0,
		ChirpCount: 0,
	}

	respondWithJSON(writer, http.StatusOK, payload)
}
