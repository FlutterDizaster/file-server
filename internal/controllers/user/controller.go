package userctrl

import (
	"context"

	"github.com/FlutterDizaster/file-server/internal/apperrors"
	jwtresolver "github.com/FlutterDizaster/file-server/internal/jwt-resolver"
	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/FlutterDizaster/file-server/internal/validator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	subject = "file-server"
)

type UserRepository interface {
	AddUser(ctx context.Context, login, passHash string) (uuid.UUID, error)
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
}

type Settings struct {
	AdminToken string
	UserRepo   UserRepository
	Resolver   *jwtresolver.JWTResolver
	Validator  *validator.Validator
}

type UserController struct {
	adminToken string
	userRepo   UserRepository
	resolver   *jwtresolver.JWTResolver
	validator  *validator.Validator
}

func New(settings Settings) *UserController {
	ctrl := &UserController{
		adminToken: settings.AdminToken,
		userRepo:   settings.UserRepo,
		resolver:   settings.Resolver,
		validator:  settings.Validator,
	}

	return ctrl
}

func (c *UserController) Register(
	ctx context.Context,
	credentials models.Credentials,
) (string, error) {
	// Verification
	if err := c.validator.ValidateCredentials(credentials); err != nil {
		return "", err
	}

	// Registration
	passHash, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	id, err := c.userRepo.AddUser(ctx, credentials.Login, string(passHash))
	if err != nil {
		return "", err
	}

	// Create token
	token, err := c.resolver.CreateToken(subject, id)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (c *UserController) Login(
	ctx context.Context,
	credentials models.Credentials,
) (string, error) {
	// Get user from the repository
	user, err := c.userRepo.GetUserByLogin(ctx, credentials.Login)
	if err != nil {
		return "", err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(credentials.Password))
	if err != nil {
		return "", apperrors.ErrWrongCredentials
	}

	// Create token
	token, err := c.resolver.CreateToken(subject, user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
