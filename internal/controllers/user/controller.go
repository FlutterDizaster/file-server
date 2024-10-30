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

// UserRepository used to get user by login and add user.
type UserRepository interface {
	// AddUser add user to repository.
	AddUser(ctx context.Context, login, passHash string) (uuid.UUID, error)

	// GetUserByLogin get user from repository.
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
}

// Settings used to create UserController.
// Settings must be provided to New function.
// All fields are required and cant be nil.
type Settings struct {
	UserRepo  UserRepository
	Resolver  *jwtresolver.JWTResolver
	Validator *validator.Validator
}

// UserController used to register and login users.
// Must be created with New function.
type UserController struct {
	userRepo  UserRepository
	resolver  *jwtresolver.JWTResolver
	validator *validator.Validator
}

// New creates new UserController.
// Returns pointer to UserController.
// Accepts Settings as argument.
func New(settings Settings) *UserController {
	ctrl := &UserController{
		userRepo:  settings.UserRepo,
		resolver:  settings.Resolver,
		validator: settings.Validator,
	}

	return ctrl
}

// Register registers new user and returns JWT token with user ID.
// Returns error if registration failed.
// Must be called with valid credentials with non-empty login, password and token.
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

// Login returns JWT token with user ID or error if login failed.
// Must be called with valid credentials with non-empty login and password.
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
