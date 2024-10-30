package handler

// FIXME: rewrite
// import (
// 	"bytes"
// 	"context"
// 	"errors"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/FlutterDizaster/file-server/internal/apperrors"
// 	"github.com/FlutterDizaster/file-server/internal/models"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func Test_userHandler(t *testing.T) {
// 	type test struct {
// 		name   string
// 		creds  models.Credentials
// 		method userCtrlMethod
// 		want   string
// 		code   int
// 	}

// 	okMethod := func(_ context.Context, _ models.Credentials) (string, error) {
// 		return "test_token", nil
// 	}

// 	internalErrorMethod := func(_ context.Context, _ models.Credentials) (string, error) {
// 		return "", errors.New("not asserted error")
// 	}

// 	wrongCredentialsMethod := func(_ context.Context, _ models.Credentials) (string, error) {
// 		return "", apperrors.ErrWrongCredentials
// 	}

// 	tests := []test{
// 		{
// 			name:   "ok test",
// 			creds:  models.Credentials{Login: "test_login", Password: "test_password"},
// 			method: okMethod,
// 			want:   "{\"response\":{\"token\":\"test_token\"}}",
// 			code:   http.StatusOK,
// 		},
// 		{
// 			name:   "internal error test",
// 			creds:  models.Credentials{Login: "test_login", Password: "test_password"},
// 			method: internalErrorMethod,
// 			want:   "{\"error\":{\"code\":500,\"text\":\"not asserted error\"}}",
// 			code:   http.StatusInternalServerError,
// 		},
// 		{
// 			name:   "wrong credentials test",
// 			creds:  models.Credentials{Login: "test_login", Password: "test_password"},
// 			method: wrongCredentialsMethod,
// 			want:   "{\"error\":{\"code\":400,\"text\":\"wrong credentials\"}}",
// 			code:   http.StatusBadRequest,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			server := httptest.NewServer(
// 				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 					userHandler(w, r, tt.method)
// 				}),
// 			)

// 			reqData, err := tt.creds.MarshalJSON()
// 			require.NoError(t, err)

// 			resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(reqData))

// 			require.NoError(t, err)
// 			defer resp.Body.Close()

// 			require.Equal(t, tt.code, resp.StatusCode)

// 			rawResp, err := io.ReadAll(resp.Body)

// 			require.NoError(t, err)

// 			assert.Equal(t, tt.want, string(rawResp))
// 		})
// 	}
// }
