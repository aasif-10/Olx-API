package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aasif-10/Olx-API/internals/httpx"
	"github.com/aasif-10/Olx-API/internals/middlewares"
	"github.com/aasif-10/Olx-API/internals/storage"
)

type UploadHandler struct {
	logger *slog.Logger
	store  *storage.Client
}

func NewUploadHandler(logger *slog.Logger, store *storage.Client) UploadHandler {
	return UploadHandler{
		logger: logger,
		store:  store,
	}
}

func (uh UploadHandler) Presign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middlewares.RequestIDFromContext(ctx)
	log := uh.logger.With("request_id", requestId)

	userId, ok := middlewares.RequireUserIdFromContext(ctx)
	if !ok {
		log.Error("no userid found in context")
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	var req PresignRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error("failed to decode", "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.CodeMalformedJSON)
		return
	}

	if len(req.Files) == 0 {
		log.Warn("file is required")
		httpx.Error(w, http.StatusBadRequest, "files must not be empty", httpx.CodeValidationFailed)
		return
	}

	if len(req.Files) > maxImagePerListing {
		log.Warn("max limit for files reached")
		httpx.Error(w, http.StatusBadRequest, fmt.Sprintf("a listing can have at most %d images", maxImagePerListing), httpx.CodeValidationFailed)
		return
	}

	for _, f := range req.Files {
		ext, ok := allowdContentTypes[f.ContentType]
		if !ok {
			log.Warn("invalid ext")
			httpx.Error(w, http.StatusBadRequest, "only image/jpeg, image/png, image/webp are allowed", httpx.CodeValidationFailed)
			return
		}

		if f.SizeBytes <= 0 || f.SizeBytes > maxImageBytes {
			log.Warn("max file size reached")
			httpx.Error(w, http.StatusBadRequest, fmt.Sprintf("size_bytes must be between 1 and %d", maxImageBytes), httpx.CodeValidationFailed)
			return
		}

		key := mintUploadKey(userId, ext)
		fmt.Print(key)

		// Actually presigning
		// uh.store.Presign
	}
	w.Write([]byte("ok"))
}
