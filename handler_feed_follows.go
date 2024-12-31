package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/FrancoMusolino/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	feedID, err := parseStringToUUID(chi.URLParam(r, "feedId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Cannot parse feed ID: %s", err))
		return
	}

	feed, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feedID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could follow the feed: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, NewApiSuccessResponse("Following feed", feed))
}

func (cfg *apiConfig) handlerUnfollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	feedID, err := parseStringToUUID(chi.URLParam(r, "feedId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Cannot parse feed ID: %s", err))
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		FeedID: feedID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could unfollow the feed: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, NewApiSuccessResponse("Unfollowing feed", struct{}{}))
}

func (cfg *apiConfig) handlerGetUserFollowedFeels(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := cfg.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not user followed feels: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, NewApiSuccessResponse(fmt.Sprintf("Found %v feeds followed by the user", len(feeds)), feeds))
}
