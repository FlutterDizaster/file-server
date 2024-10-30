package application

import (
	"context"
	"net/http"
	"time"

	docctrl "github.com/FlutterDizaster/file-server/internal/controllers/document"
	userctrl "github.com/FlutterDizaster/file-server/internal/controllers/user"
	jwtresolver "github.com/FlutterDizaster/file-server/internal/jwt-resolver"
	"github.com/FlutterDizaster/file-server/internal/migrator"
	"github.com/FlutterDizaster/file-server/internal/repository/miniorepo"
	"github.com/FlutterDizaster/file-server/internal/repository/postgresrepo"
	"github.com/FlutterDizaster/file-server/internal/repository/redisrepo"
	"github.com/FlutterDizaster/file-server/internal/server"
	"github.com/FlutterDizaster/file-server/internal/server/handler"
	"github.com/FlutterDizaster/file-server/internal/validator"
	"github.com/FlutterDizaster/file-server/pkg/configloader"
)

const (
	// shutdownMaxTime is the maximum time allowed to gracefully shutdown the server.
	shutdownMaxTime = 5 * time.Second
)

// Service represents the application service.
type Service interface {
	Start(ctx context.Context) error
}

//nolint:lll // struct tags too long
type Settings struct {
	PostgresConnectionString string `desc:"postgres connection string" env:"DATABASE_DSN"       name:"database-dsn"       short:"d"`
	PostgresMigrationsPath   string `desc:"postgres migrations path"   env:"DB_MIGRATIONS_PATH" name:"db-migrations-path" short:"m"`

	RedisConnectionString string `desc:"redis connection string"      env:"REDIS_DSN"       name:"redis-DSN"       short:"r"`
	RedisCacheTTL         string `desc:"redis cache ttl, default 24h" env:"REDIS_CACHE_TTL" name:"redis-cache-ttl" short:"t" default:"24h"`

	MinioEndpoint  string `desc:"minio endpoint"   env:"MINIO_ENDPOINT"   name:"minio-endpoint"   short:"e"`
	MinioAccessKey string `desc:"minio access key" env:"MINIO_ACCESS_KEY" name:"minio-access-key" short:"a"`
	MinioSecretKey string `desc:"minio secret key" env:"MINIO_SECRET_KEY" name:"minio-secret-key" short:"s"`
	MinioBucket    string `desc:"minio bucket"     env:"MINIO_BUCKET"     name:"minio-bucket"     short:"b"`
	MinioUseSSL    bool   `desc:"minio use ssl"    env:"MINIO_USE_SSL"    name:"minio-use-ssl"    short:"u"`

	AdminToken string `desc:"admin token" env:"ADMIN_TOKEN" name:"admin-token"`

	JWTSecret string `desc:"jwt secret"                      env:"JWT_SECRET" name:"jwt-secret" short:"j"`
	JWTIssuer string `desc:"jwt issuer, default file-server" env:"JWT_ISSUER" name:"jwt-issuer"           default:"file-server"`
	JWTTTL    string `desc:"jwt ttl, default 24h"            env:"JWT_TTL"    name:"jwt-ttl"              default:"24h"`

	HTTPAddr                 string `desc:"http address, default localhost"             env:"HTTP_ADDR"            name:"http-addr"            short:"a" default:"localhost"`
	HTTPPort                 string `desc:"http port, default 8080"                     env:"HTTP_PORT"            name:"http-port"            short:"p" default:"8080"`
	HandlerMaxUploadFileSize int64  `desc:"handler max upload file size, default 200Mb" env:"MAX_UPLOAD_FILE_SIZE" name:"max-upload-file-size"           default:"209715200"`
}

// New creates a new application service that can be started with Start method.
func New(ctx context.Context) (Service, error) {
	// Load config
	var settings Settings
	err := configloader.LoadConfig(&settings)
	if err != nil {
		return nil, err
	}

	// Run migrations
	err = migrator.RunMigrations(
		ctx,
		settings.PostgresConnectionString,
		settings.PostgresMigrationsPath,
	)
	if err != nil {
		return nil, err
	}

	// new repositories
	postgresRepo, err := newPostgresRepository(ctx, settings)
	if err != nil {
		return nil, err
	}

	redisRepo, err := newRedisRepository(ctx, settings)
	if err != nil {
		return nil, err
	}

	minioRepo, err := newMinioRepository(ctx, settings)
	if err != nil {
		return nil, err
	}

	// new resolver and validator
	resolver, err := newJWTResolver(settings)
	if err != nil {
		return nil, err
	}

	validator, err := newValidator(settings)
	if err != nil {
		return nil, err
	}

	// new controllers
	documentsController := newDocumentsController(
		minioRepo,
		postgresRepo,
		postgresRepo,
		redisRepo,
	)

	userController := newUserController(
		postgresRepo,
		resolver,
		validator,
	)

	// new Handler
	handler := newHandler(
		resolver,
		userController,
		documentsController,
		settings.HandlerMaxUploadFileSize,
	)

	// new server
	server := newServer(settings, handler)

	return server, nil
}

func newPostgresRepository(
	ctx context.Context,
	settings Settings,
) (*postgresrepo.PostgresRepository, error) {
	return postgresrepo.New(ctx, settings.PostgresConnectionString)
}

func newRedisRepository(
	ctx context.Context,
	settings Settings,
) (*redisrepo.RedisRepository, error) {
	ttl, err := time.ParseDuration(settings.RedisCacheTTL)
	if err != nil {
		return nil, err
	}
	repoSettings := redisrepo.Settings{
		ConnectionString: settings.RedisConnectionString,
		TTL:              ttl,
	}

	return redisrepo.New(ctx, repoSettings)
}

func newMinioRepository(
	ctx context.Context,
	settings Settings,
) (*miniorepo.MinioRepository, error) {
	repoSettings := miniorepo.Settings{
		Endpoint:  settings.MinioEndpoint,
		AccessKey: settings.MinioAccessKey,
		SecretKey: settings.MinioSecretKey,
		Bucket:    settings.MinioBucket,
		UseSSL:    settings.MinioUseSSL,
	}

	return miniorepo.New(ctx, repoSettings)
}

func newJWTResolver(settings Settings) (*jwtresolver.JWTResolver, error) {
	ttl, err := time.ParseDuration(settings.JWTTTL)
	if err != nil {
		return nil, err
	}
	jwtSettings := jwtresolver.Settings{
		Secret:   settings.JWTSecret,
		Issuer:   settings.JWTIssuer,
		TokenTTL: ttl,
	}

	return jwtresolver.New(jwtSettings), nil
}

func newValidator(settings Settings) (*validator.Validator, error) {
	return validator.New(settings.AdminToken)
}

func newDocumentsController(
	fileRepo docctrl.FileRepository,
	userRepo docctrl.UserRepository,
	metaRepo docctrl.MetadataRepository,
	cache docctrl.MetadataCache,
) *docctrl.DocumentsController {
	controllerSettings := docctrl.Settings{
		FileRepo: fileRepo,
		MetaRepo: metaRepo,
		UserRepo: userRepo,
		Cache:    cache,
	}

	return docctrl.New(controllerSettings)
}

func newUserController(
	userRepo userctrl.UserRepository,
	resolver *jwtresolver.JWTResolver,
	validator *validator.Validator,
) *userctrl.UserController {
	controllerSettings := userctrl.Settings{
		UserRepo:  userRepo,
		Resolver:  resolver,
		Validator: validator,
	}

	return userctrl.New(controllerSettings)
}

func newHandler(
	resolver *jwtresolver.JWTResolver,
	userCtrl handler.UserController,
	docCtrl handler.DocumentsController,
	maxUploadSize int64,
) *handler.Handler {
	handlerSettings := handler.Settings{
		JWTResolver:       resolver,
		UserCtrl:          userCtrl,
		DocumentsCtrl:     docCtrl,
		MaxUploadFileSize: maxUploadSize,
	}

	return handler.New(handlerSettings)
}

func newServer(settings Settings, handler http.Handler) *server.Server {
	serverSettings := server.Settings{
		Addr:    settings.HTTPAddr,
		Port:    settings.HTTPPort,
		Handler: handler,

		ShutdownMaxTime: shutdownMaxTime,
	}

	return server.New(serverSettings)
}
