package api

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/FlutterDizaster/file-server/internal/api/middlewares"
	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/google/uuid"
)

func (a *API) docGetListHandler(w http.ResponseWriter, r *http.Request) {
	respData := a.prepareDocListResponse(w, r)
	if respData == nil {
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(respData); err != nil {
		slog.Error("Error while writing response", slog.Any("err", err))
		return
	}
}

func (a *API) docGetListHeadHandler(w http.ResponseWriter, r *http.Request) {
	respData := a.prepareDocListResponse(w, r)
	if respData == nil {
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Lenght", strconv.Itoa(len(respData)))
	w.WriteHeader(http.StatusOK)
}

func (a API) prepareDocListResponse(
	w http.ResponseWriter,
	r *http.Request,
) []byte {
	userID, ok := r.Context().Value(middlewares.KeyUserID).(uuid.UUID)
	if !ok {
		slog.Error("User id not found in context")
		responseWithError(w, r, nil, "User id not found")
		return nil
	}

	// Check content type
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		responseWithError(w, r, nil, "Invalid content type")
		return nil
	}

	// Reading body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responseWithError(w, r, err, "Error while reading body")
		return nil
	}
	defer r.Body.Close()

	var filesListReq models.FilesListRequest
	if err = filesListReq.UnmarshalJSON(body); err != nil {
		responseWithError(w, r, err, "Error while unmarshaling body")
		return nil
	}

	// Execute method
	files, err := a.documentsCtrl.GetFilesInfo(r.Context(), userID, filesListReq)
	if err != nil {
		responseWithError(w, r, err, "Error while getting files list")
		return nil
	}

	// Prepare response
	resp := &models.Response{
		Data: &models.ResponseFilesList{
			Docs: files,
		},
	}

	// Marshalling response
	respData, err := resp.MarshalJSON()
	if err != nil {
		responseWithError(w, r, err, "Error while marshaling response")
		return nil
	}

	return respData
}
