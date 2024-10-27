package api

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/FlutterDizaster/file-server/internal/api/middlewares"
	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/google/uuid"
)

func (a API) docPostHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.KeyUserID).(uuid.UUID)
	if !ok {
		slog.Error("User id not found in context")
		responseWithError(w, r, nil, "User id not found")
		return
	}
	// Check content-type
	if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		responseWithError(w, r, nil, "Wrong content type")
		return
	}

	// Parsing multipart form
	err := r.ParseMultipartForm(a.maxUploadFileSize)
	if err != nil {
		responseWithError(w, r, err, "Error while parsing multipart form")
		return
	}

	// Extract metadata
	metaStr := r.FormValue("meta")
	var metadata models.Metadata
	if err = metadata.UnmarshalJSON([]byte(metaStr)); err != nil {
		responseWithError(w, r, err, "Error while unmarshaling metadata")
		return
	}

	// Set file owner ID
	metadata.OwnerID = &userID

	// Extract JSON data
	jsonStr := r.FormValue("json")
	metadata.JSON = models.JSONString(jsonStr)

	// Extract file
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		responseWithError(w, r, err, "Error while getting file")
		return
	}
	defer file.Close()

	// Set file size
	metadata.FileSize = fileHeader.Size

	// Upload document
	if err = a.documentsCtrl.UploadDocument(r.Context(), metadata, file); err != nil {
		responseWithError(w, r, err, "Error while uploading document")
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
		responseWithError(w, r, err, "Error while marshaling response")
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
