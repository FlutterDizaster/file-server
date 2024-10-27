package api

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/FlutterDizaster/file-server/internal/models"
)

type userCtrlMethod func(ctx context.Context, cred models.Credentials) (string, error)

func userHandler(w http.ResponseWriter, r *http.Request, method userCtrlMethod) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		responseWithError(w, r, nil, "wrong content type")
		return
	}

	// Extract credentials
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responseWithError(w, r, err, "Extracting credentials failed")
		return
	}
	defer r.Body.Close()

	var cred models.Credentials
	if err = cred.UnmarshalJSON(body); err != nil {
		responseWithError(w, r, err, "Error while unmarshaling credentials")
		return
	}

	// Execute method
	token, err := method(r.Context(), cred)
	if err != nil {
		responseWithError(w, r, err, "User login/registration failed")
		return
	}

	// Create response
	resp := models.Response{
		Response: &models.Credentials{
			Token: token,
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
	_, err = w.Write(respData)
	if err != nil {
		slog.Error("Error while writing response", slog.Any("err", err))
	}
}
