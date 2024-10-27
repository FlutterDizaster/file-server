package api

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	jwtresolver "github.com/FlutterDizaster/file-server/internal/jwt-resolver"
	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// UserController used to register and login users.
type UserController interface {
	// Returns jwt token with user ID or error if user creation failed.
	// Must be called with valid credentials with non-empty login, password and token.
	Register(ctx context.Context, credentials models.Credentials) (string, error)

	// Login returns jwt token with user ID or error if login failed.
	// Must be called with valid credentials with non-empty login and password.
	Login(ctx context.Context, credentials models.Credentials) (string, error)
}

type DocumentsController interface {
	UploadDocument(ctx context.Context, meta models.Metadata, file io.Reader) error
	GetFilesInfo(
		ctx context.Context,
		userID uuid.UUID,
		filesListRequest models.FilesListRequest,
	) ([]models.Metadata, error)
	GetFileInfo(ctx context.Context, id, userID uuid.UUID) (*models.Metadata, error)
	GetFile(Ctx context.Context, id uuid.UUID) (io.ReadSeeker, error)
	DeleteFile(ctx context.Context, id, userID uuid.UUID) error
}

type Settings struct {
	Addr          string
	Port          string
	JWTResolver   *jwtresolver.JWTResolver
	UserCtrl      UserController
	DocumentsCtrl DocumentsController

	ShutdownMaxTime   time.Duration
	MaxUploadFileSize int64
}

type API struct {
	server            *http.Server
	addr              string
	port              string
	jwtResolver       *jwtresolver.JWTResolver
	userCtrl          UserController
	documentsCtrl     DocumentsController
	shutdownMaxTime   time.Duration
	maxUploadFileSize int64
}

func New(settings Settings) *API {
	a := &API{
		addr:          settings.Addr,
		port:          settings.Port,
		jwtResolver:   settings.JWTResolver,
		userCtrl:      settings.UserCtrl,
		documentsCtrl: settings.DocumentsCtrl,

		shutdownMaxTime:   settings.ShutdownMaxTime,
		maxUploadFileSize: settings.MaxUploadFileSize,
	}

	router := a.setupRouter()

	a.server = &http.Server{
		ReadHeaderTimeout: time.Second,
		Addr:              fmt.Sprintf("%s:%s", a.addr, a.port),
		Handler:           router,
	}

	return a
}

func (a *API) Start(ctx context.Context) error {
	eg, egCtx := errgroup.WithContext(ctx)

	errorOnStart := false

	eg.Go(func() error {
		<-egCtx.Done()

		// Skip shutdown if error on start
		if errorOnStart {
			return nil
		}

		slog.Info("Shutting down server")
		shutdownCtx, cancle := context.WithTimeout(context.Background(), a.shutdownMaxTime)
		defer cancle()

		return a.server.Shutdown(shutdownCtx)
	})

	eg.Go(func() error {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Error while starting server", slog.Any("err", err))
			errorOnStart = true
			return err
		}
		return nil
	})

	return eg.Wait()
}
