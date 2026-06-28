package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"so-many-v2/realtime_comments/pkg/http_tools"

	"github.com/go-chi/chi/v5"
)

func (rh *RouterHandler) SSEHandler(w http.ResponseWriter, req *http.Request) {
	postIDStr := chi.URLParam(req, "post_id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		http_tools.Error(w, http.StatusUnprocessableEntity, errors.New("invalid post id"))
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http_tools.Error(w, http.StatusInternalServerError, errors.New("streaming unsupported"))
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	
	channel := fmt.Sprintf("post:%d", postID)
	sub := rh.service.Connections.Add(channel)
	defer rh.service.Connections.Remove(channel, sub)

	flusher.Flush()

	for {
		select {
		case <-req.Context().Done():
			return
		case <-rh.service.Connections.Done():
			return
		case msg := <-sub.Send():
			if _, err := fmt.Fprintf(w, "data: %s\n\n", msg); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}