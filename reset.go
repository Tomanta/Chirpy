package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

func (cfg *apiConfig) handlerReset(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("DEBUG: Resetting...")

	platform := os.Getenv("PLATFORM")

	if platform != "dev" {
		respondWithError(writer, http.StatusForbidden, "", nil)
		return
	}

	type Reset struct {
		HitCount  int `json:"hit_count"`
		UserCount int `json:"user_count"`
	}

	cfg.fileserverHits.Store(0)

	err := cfg.dbQueries.ResetUsers(context.Background())
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Couldn't reset users", err)
		return
	}

	payload := Reset{
		HitCount:  0,
		UserCount: 0,
	}

	respondWithJSON(writer, http.StatusOK, payload)
}
