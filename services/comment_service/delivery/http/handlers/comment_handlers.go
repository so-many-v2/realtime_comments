package handlers

import (
	"errors"
	"net/http"
	"so-many-v2/realtime_comments/pkg/errors/agg_errors"
	"so-many-v2/realtime_comments/pkg/errors/db_errors"
	"so-many-v2/realtime_comments/pkg/errors/http_errors"
	"so-many-v2/realtime_comments/pkg/http_tools"
	"so-many-v2/realtime_comments/services/models"
	"strconv"
	"time"
)

func (rh *RouterHandler) GetPostHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	postIDStr := req.URL.Query().Get("post_id")
	timeFromStr := req.URL.Query().Get("time_from")
	limitStr := req.URL.Query().Get("limit")

	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		http_tools.Error(w, http.StatusUnprocessableEntity, errors.New("invalid post_id"))
		return
	}

	sec, err := strconv.ParseInt(timeFromStr, 10, 64)
	if err != nil {
		http_tools.Error(w, http.StatusUnprocessableEntity, errors.New("invalid time_from"))
		return
	}
	timeFrom := time.Unix(sec, 0).UTC()

	limit, err := strconv.ParseUint(limitStr, 10, 64)
	if err != nil {
		http_tools.Error(w, http.StatusUnprocessableEntity, errors.New("invalid limit"))
		return
	}

	comments, err := rh.service.Comments.GetPostCommentBatch(ctx, uint(postID), timeFrom, uint(limit))
	if err != nil {
		var ve agg_errors.ValidationError
		if errors.As(err, &ve) {
			http_tools.Error(w, http.StatusUnprocessableEntity, ve)
			return
		}
		if errors.Is(err, db_errors.ErrNotFound) {
			http_tools.Error(w, http.StatusNotFound, errors.New("post not found"))
			return
		}
		http_tools.Error(w, http.StatusInternalServerError, http_errors.ErrServerError)
		return
	}
	http_tools.JSON(w, http.StatusOK, comments)
}

func (rh *RouterHandler) CreatePostHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var data models.CreateComment
	decoded := http_tools.NewDecoder(w, req)
	if err := decoded.Decode(&data); err != nil {
		http_tools.Error(w, http.StatusUnprocessableEntity, errors.New("invalid post data"))
		return
	}

	postId, err := rh.service.Comments.CreateComment(ctx, &data)
	if err != nil {
		var ve agg_errors.ValidationError
		if errors.As(err, &ve) {
			http_tools.Error(w, http.StatusUnprocessableEntity, ve)
			return
		}
		http_tools.Error(w, http.StatusInternalServerError, http_errors.ErrServerError)
		return
	}
	http_tools.JSON(w, http.StatusCreated, map[string]uint{
		"comment_id": postId,
	})
}
