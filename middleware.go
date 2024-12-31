package main

import (
	"fmt"
	"net/http"

	"github.com/FrancoMusolino/rssagg/internal/auth"
	"github.com/FrancoMusolino/rssagg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't get the user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
