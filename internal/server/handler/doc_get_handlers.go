package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/FlutterDizaster/file-server/internal/server/middlewares"
	"github.com/google/uuid"
)

type serveFileStrategy func(http.ResponseWriter, *http.Request, models.Metadata)

func (h Handler) docGetHandler(w http.ResponseWriter, r *http.Request) {
	strategyMap := map[bool]serveFileStrategy{
		true:  h.serveBinaryFileHandler,
		false: h.serveJSONFileHandler,
	}

	// Get user id
	userID, ok := r.Context().Value(middlewares.KeyUserID).(uuid.UUID)
	if !ok {
		h.responseWithError(w, r, nil, "User id not found")
		return
	}

	// Get doc id
	docID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.responseWithError(w, r, err, "Invalid document id")
		return
	}

	// Get file info
	info, err := h.documentsCtrl.GetFileInfo(r.Context(), docID, userID)
	if err != nil {
		h.responseWithError(w, r, err, "Error while getting file info")
		return
	}

	// Send response
	strategyMap[info.File](w, r, info)
}

func (h Handler) docGetHeadHandler(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userID, ok := r.Context().Value(middlewares.KeyUserID).(uuid.UUID)
	if !ok {
		h.responseWithError(w, r, nil, "User id not found")
		return
	}

	// Get doc id
	docID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.responseWithError(w, r, err, "Invalid document id")
		return
	}

	// Get file info
	info, err := h.documentsCtrl.GetFileInfo(r.Context(), docID, userID)
	if err != nil {
		h.responseWithError(w, r, err, "Error while getting file info")
		return
	}

	// Send response
	if info.File {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+info.Name)
		w.Header().Set("Content-Lenght", strconv.FormatInt(info.FileSize, 10))
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Lenght", strconv.Itoa(len(info.JSON)))

	w.WriteHeader(http.StatusOK)
}

func (h Handler) serveBinaryFileHandler(
	w http.ResponseWriter,
	r *http.Request,
	meta models.Metadata,
) {
	// Get file
	file, err := h.documentsCtrl.GetFile(r.Context(), meta)
	if err != nil {
		h.responseWithError(w, r, err, "Error while getting file")
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+meta.Name)
	w.Header().Set("Content-Lenght", strconv.FormatInt(meta.FileSize, 10))

	http.ServeContent(w, r, meta.Name, time.Now(), file)
}

func (h Handler) serveJSONFileHandler(
	w http.ResponseWriter,
	r *http.Request,
	meta models.Metadata,
) {
	respData := models.Response{
		Data: &meta.JSON,
	}

	resp, err := respData.MarshalJSON()
	if err != nil {
		h.responseWithError(w, r, err, "Error while marshaling response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(resp)
	if err != nil {
		slog.Error("Error while writing response", slog.Any("err", err))
		return
	}
}
