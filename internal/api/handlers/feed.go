package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"je-suis-ici-activitypub/internal/services"
	"net/http"
	"strconv"
)

type FeedHandler struct {
	checkinService services.CheckinService
}

func NewFeedHandler(checkinService services.CheckinService) *FeedHandler {
	return &FeedHandler{
		checkinService: checkinService,
	}
}

func (fh *FeedHandler) RegisterFeedRouters(r chi.Router) {
	r.Get("/feed", fh.GetGlobalFeed)
}

func (fh *FeedHandler) GetGlobalFeed(w http.ResponseWriter, r *http.Request) {
	// get pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// get global feed
	checkins, err := fh.checkinService.GetGlobalFeed(r.Context(), page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return global feed
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"checkins":  checkins,
		"page":      page,
		"page_size": pageSize,
	})
}
