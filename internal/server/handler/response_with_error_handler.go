package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/FlutterDizaster/file-server/internal/apperrors"
	"github.com/FlutterDizaster/file-server/internal/models"
)

func (h Handler) responseWithError(w http.ResponseWriter, r *http.Request, err error, msg string) {
	resp := &models.Response{
		Error: &models.ResponseError{
			Text: msg,
		},
	}
	var appserror *apperrors.Error

	switch {
	case errors.As(err, &appserror):
		resp.Error.Code = appserror.Code
		resp.Error.Text = fmt.Sprintf("%s: %s", msg, appserror.Message)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Error.Code)

	if _, err = w.Write(respData); err != nil {
		slog.Error("Error while writing response", slog.Any("err", err))
		return
	}
}
