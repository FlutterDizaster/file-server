package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/FlutterDizaster/file-server/internal/apperrors"
	"github.com/FlutterDizaster/file-server/internal/models"
)

func responseWithError(w http.ResponseWriter, r *http.Request, err error, msg string) {
	resp := &models.Response{
		Error: &models.ResponseError{
			Text: msg,
		},
	}

	switch e := err.(type) {
	case *apperrors.Error:
		resp.Error.Code = e.Code
		resp.Error.Text = fmt.Sprintf("%s: %s", msg, e.Message)
	case nil:
		resp.Error.Code = http.StatusBadRequest
		resp.Error.Text = msg
	default:
		slog.Error(
			"Error while processing request",
			slog.String("Message", msg),
			slog.String("Method", r.Method),
			slog.String("URL", r.URL.String()),
			slog.Any("err", err),
		)
		resp.Error.Code = http.StatusInternalServerError
		resp.Error.Text = err.Error()
	}

	respData, err := resp.MarshalJSON()
	if err != nil {
		slog.Error("Error while marshaling response", slog.Any("err", err))
		return
	}

	w.WriteHeader(resp.Error.Code)
	w.Header().Set("Content-Type", "application/json")

	if _, err = w.Write(respData); err != nil {
		slog.Error("Error while writing response", slog.Any("err", err))
		return
	}
}
