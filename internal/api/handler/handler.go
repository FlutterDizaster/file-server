package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/FlutterDizaster/file-server/internal/api/middlewares"
	jwtresolver "github.com/FlutterDizaster/file-server/internal/jwt-resolver"
	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/google/uuid"
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
	GetFileInfo(ctx context.Context, id, userID uuid.UUID) (models.Metadata, error)
	GetFile(Ctx context.Context, id uuid.UUID) (io.ReadSeeker, error)
	DeleteFile(ctx context.Context, id, userID uuid.UUID) error
}

type Settings struct {
	JWTResolver       *jwtresolver.JWTResolver
	UserCtrl          UserController
	DocumentsCtrl     DocumentsController
	MaxUploadFileSize int64
}

type Handler struct {
	router            *http.ServeMux
	jwtResolver       *jwtresolver.JWTResolver
	userCtrl          UserController
	documentsCtrl     DocumentsController
	maxUploadFileSize int64
}

func New(settings Settings) *Handler {
	h := &Handler{
		jwtResolver:       settings.JWTResolver,
		userCtrl:          settings.UserCtrl,
		documentsCtrl:     settings.DocumentsCtrl,
		maxUploadFileSize: settings.MaxUploadFileSize,
	}

	h.setupRouter()

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// setupRouter sets up the router with the appropriate middleware chains
// and returns the router.
func (h *Handler) setupRouter() {
	router := http.NewServeMux()

	// Public routes
	userRouter := http.NewServeMux()
	userRouter.HandleFunc("POST /auth", h.userAuthHandler)
	userRouter.HandleFunc("POST /register", h.userRegisterHandler)

	// Private routes
	docRouter := http.NewServeMux()
	docRouter.HandleFunc("GET /{id}", h.docGetHandler)
	docRouter.HandleFunc("HEAD /{id}", h.docGetHeadHandler)
	docRouter.HandleFunc("GET /", h.docGetListHandler)
	docRouter.HandleFunc("HEAD /", h.docGetListHeadHandler)
	docRouter.HandleFunc("POST /", h.docPostHandler)
	docRouter.HandleFunc("DELETE /{id}", h.docDeleteHandler)

	// Public middleware chain
	publicChain := middlewares.MakeChain(
		middlewares.Logger,
	)

	// Private middleware chain
	authMw := middlewares.Auth{
		Resolver: h.jwtResolver,
	}
	privateChain := middlewares.MakeChain(
		middlewares.Logger,
		authMw.Handle,
	)

	// Setup general router
	router.Handle("/api/", publicChain(http.StripPrefix("/api/", userRouter)))
	router.Handle("/api/docs", privateChain(http.StripPrefix("/api/docs", docRouter)))

	h.router = router
}
