package main

import (
	"net/http"
	"strconv"

	"github.com/chiagxziem/gater/internal/cursor"
	"github.com/chiagxziem/gater/internal/jsonutil"
	"github.com/chiagxziem/gater/internal/store"
)

func (a *application) getPublishedEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := loggerFromCtx(ctx)

	limit := 50
	if lq := r.URL.Query().Get("limit"); lq != "" {
		if v, err := strconv.Atoi(lq); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	var cur *cursor.Cursor
	if cq := r.URL.Query().Get("cursor"); cq != "" {
		c, err := cursor.Decode(cq)
		if err != nil {
			jsonutil.WriteError(w, http.StatusBadRequest, "invalid cursor")
			return
		}
		cur = &c
	}

	// get events up to limit + 1 to ensure there's a next page
	events, err := a.store.Events.GetPublished(ctx, cur, limit+1)
	if err != nil {
		logger.Error("failed to get published events", "error", err)
		jsonutil.WriteError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	//  if there's a next page, create next cursor
	var nextCursor string
	if len(events) == limit+1 {
		events = events[:limit]
		last := events[len(events)-1]
		c := cursor.Cursor{
			Timestamp: last.StartsAt,
			ID:        last.ID,
		}
		nextCursor = cursor.Encode(c)
	}

	type returnData struct {
		Message    string         `json:"message"`
		Events     []*store.Event `json:"events"`
		NextCursor string         `json:"next_cursor"`
	}
	jsonutil.WriteData(w, http.StatusOK, returnData{
		Message:    "published events retrieved successfully",
		Events:     events,
		NextCursor: nextCursor,
	})
}
