package api

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/FlutterDizaster/file-server/internal/api/middlewares"
	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/google/uuid"
)

func (a API) docGetHandler(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userID, ok := r.Context().Value(middlewares.KeyUserID).(uuid.UUID)
	if !ok {
		responseWithError(w, r, nil, "User id not found")
		return
	}

	// Get doc id
	docID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		responseWithError(w, r, err, "Invalid document id")
		return
	}

	// Get file info
	info, err := a.documentsCtrl.GetFileInfo(r.Context(), docID, userID)
	if err != nil {
		responseWithError(w, r, err, "Error while getting file info")
		return
	}

	// Send response
	// If requested document is file - send file
	// Othervise send JSON
	if info.File {
		// Get file
		var file io.ReadSeeker
		file, err = a.documentsCtrl.GetFile(r.Context(), docID)
		if err != nil {
			responseWithError(w, r, err, "Error while getting file")
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+info.Name)
		w.Header().Set("Content-Lenght", strconv.FormatInt(info.FileSize, 10))

		http.ServeContent(w, r, info.Name, time.Now(), file)
	} else {
		respData := models.Response{
			Data: &info.JSON,
		}

		var resp []byte
		resp, err = respData.MarshalJSON()
		if err != nil {
			responseWithError(w, r, err, "Error while marshaling response")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(resp)
		if err != nil {
			slog.Error("Error while writing response", slog.Any("err", err))
			return
		}
	}
}

func (a API) docGetHeadHandler(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userID, ok := r.Context().Value(middlewares.KeyUserID).(uuid.UUID)
	if !ok {
		responseWithError(w, r, nil, "User id not found")
		return
	}

	// Get doc id
	docID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		responseWithError(w, r, err, "Invalid document id")
		return
	}

	// Get file info
	info, err := a.documentsCtrl.GetFileInfo(r.Context(), docID, userID)
	if err != nil {
		responseWithError(w, r, err, "Error while getting file info")
		return
	}

	// Send response
	if info.File {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+info.Name)
		w.Header().Set("Content-Lenght", strconv.FormatInt(info.FileSize, 10))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Lenght", strconv.Itoa(len(info.JSON)))
	}

	w.WriteHeader(http.StatusOK)
}
