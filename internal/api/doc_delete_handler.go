package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/FlutterDizaster/file-server/internal/api/middlewares"
	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/google/uuid"
)

func (a *API) docDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userID, ok := r.Context().Value(middlewares.KeyUserID).(uuid.UUID)
	if !ok {
		slog.Error("User id not found in context")
		responseWithError(w, r, nil, "User id not found")
		return
	}

	// Get doc id
	docID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		responseWithError(w, r, err, "Invalid document id")
		return
	}

	// Delete document
	err = a.documentsCtrl.DeleteFile(r.Context(), docID, userID)
	if err != nil {
		responseWithError(w, r, err, "Error while deleting document")
		return
	}

	// Prepare response
	respString := models.JSONString(fmt.Sprintf("{%s: true}", docID))
	respData := models.Response{
		Response: &respString,
	}

	// Marshal response
	resp, err := respData.MarshalJSON()
	if err != nil {
		slog.Error("Error while marshaling response", slog.Any("err", err))
		responseWithError(w, r, err, "Error while marshaling response")
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(resp); err != nil {
		slog.Error("Error while writing response", slog.Any("err", err))
		return
	}
}
