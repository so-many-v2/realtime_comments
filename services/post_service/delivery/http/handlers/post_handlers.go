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

	"github.com/go-chi/chi/v5"
)

func (rh *RouterHandler) GetPostHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	postIdStr := chi.URLParam(req, "postID")
	postID, err := strconv.ParseUint(postIdStr, 10, 64)
	if err != nil {
		http_tools.Error(w, http.StatusUnprocessableEntity, errors.New("invalid post id"))
		return
	}

	post, err := rh.service.Post.GetPost(ctx, uint(postID))
	if err != nil {
		if errors.Is(err, db_errors.ErrNotFound) {
			http_tools.Error(w, http.StatusNotFound, errors.New("post not found"))
			return
		}
		http_tools.Error(w, http.StatusNotFound, http_errors.ErrServerError)
		return
	}
	http_tools.JSON(w, http.StatusOK, post)
}

func (rh *RouterHandler) ListPostsHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var limit uint = 20
	if raw := req.URL.Query().Get("limit"); raw != "" {
		if v, err := strconv.ParseUint(raw, 10, 64); err == nil {
			limit = uint(v)
		}
	}

	var offset uint = 0
	if raw := req.URL.Query().Get("offset"); raw != "" {
		if v, err := strconv.ParseUint(raw, 10, 64); err == nil {
			offset = uint(v)
		}
	}

	posts, err := rh.service.Post.ListPosts(ctx, limit, offset)
	if err != nil {
		http_tools.Error(w, http.StatusInternalServerError, http_errors.ErrServerError)
		return
	}
	http_tools.JSON(w, http.StatusOK, posts)
}

func (rh *RouterHandler) CreatePostHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var postData models.CreatePost
	decoded := http_tools.NewDecoder(w, req)
	if err := decoded.Decode(&postData); err != nil {
		http_tools.Error(w, http.StatusUnprocessableEntity, errors.New("invalid post data"))
		return
	}

	postId, err := rh.service.Post.CreatePost(ctx, &postData)
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
		"post_id": postId,
	})
}
