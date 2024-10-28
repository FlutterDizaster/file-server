package handler

import (
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/FlutterDizaster/file-server/internal/apperrors"
	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/FlutterDizaster/file-server/internal/server/middlewares"
	"github.com/google/uuid"
)

func (h Handler) docPostHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.KeyUserID).(uuid.UUID)
	if !ok {
		slog.Error("User id not found in context")
		h.responseWithError(w, r, nil, "User id not found")
		return
	}
	// Check content-type
	if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		err := apperrors.ErrInvalidContentType
		h.responseWithError(w, r, err, r.Header.Get("Content-Type"))
		return
	}

	// Parsing multipart form
	err := r.ParseMultipartForm(h.maxUploadFileSize)
	if err != nil {
		h.responseWithError(w, r, err, "Error while parsing multipart form")
		return
	}

	// Extract metadata
	metaStr := r.FormValue("meta")
	var metadata models.Metadata
	if err = metadata.UnmarshalJSON([]byte(metaStr)); err != nil {
		h.responseWithError(w, r, err, "Error while unmarshaling metadata")
		return
	}

	// Set file owner ID
	metadata.OwnerID = &userID

	// Extract JSON data
	jsonStr := r.FormValue("json")
	metadata.JSON = models.JSONString(jsonStr)

	// Extract file
	var file io.ReadCloser
	if metadata.File {
		var fileHeader *multipart.FileHeader
		file, fileHeader, err = r.FormFile("file")
		if err != nil {
			h.responseWithError(w, r, err, "Error while getting file")
			return
		}
		defer file.Close()

		// Set file size
		metadata.FileSize = fileHeader.Size
	}

	// Upload document
	if err = h.documentsCtrl.UploadDocument(r.Context(), metadata, file); err != nil {
		h.responseWithError(w, r, err, "Error while uploading document")
		return
	}

	// Prepare response
	resp := models.Response{
		Data: &models.ResponseUploading{
			JSON: models.JSONString(jsonStr),
			File: metadata.Name,
		},
	}

	// Marshal response
	respData, err := resp.MarshalJSON()
	if err != nil {
		h.responseWithError(w, r, err, "Error while marshaling response")
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(respData); err != nil {
		slog.Error("Error while writing response", slog.Any("err", err))
		return
	}
}
