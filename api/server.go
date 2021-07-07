package main

import (
	"context"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shellhub-io/shellhub/api/apicontext"
	"github.com/shellhub-io/shellhub/api/routes"
	"github.com/shellhub-io/shellhub/api/routes/middlewares"
	storecache "github.com/shellhub-io/shellhub/api/store/cache"
	"github.com/shellhub-io/shellhub/api/store/mongo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var serverCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return startServer()
	},
}

// Provides the configuration for the API service.
// The values are load from the system environment variables.
type config struct {
	// MongoDB connection string (URI format)
	MongoURI string `envconfig:"mongo_uri" default:"mongodb://mongo:27017"`
	// Redis connection stirng (URI format)
	RedisURI string `envconfig:"redis_uri" default:"redis://redis:6379"`
	// Enable store cache
	StoreCache bool `envconfig:"store_cache" default:"false"`
}

func startServer() error {
	logrus.Info("Starting API server")

	e := echo.New()
	e.Use(middleware.Logger())

	// Populates configuration based on environment variables prefixed with 'API_'
	var cfg config
	if err := envconfig.Process("api", &cfg); err != nil {
		logrus.WithError(err).Fatal("Failed to load environment variables")
	}

	logrus.Info("Connecting to MongoDB")

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongodriver.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to MongoDB")
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		logrus.WithError(err).Fatal("Failed to ping MongoDB")
	}

	logrus.Info("Running database migrations")

	if err := mongo.ApplyMigrations(client.Database("main")); err != nil {
		logrus.WithError(err).Fatal("Failed to apply mongo migrations")
	}

	var cache storecache.Cache
	if cfg.StoreCache {
		logrus.Info("Using redis as store cache backend")

		cache, err = storecache.NewRedisCache(cfg.RedisURI)
		if err != nil {
			logrus.WithError(err).Error("Failed to configure redis store cache")
		}
	} else {
		logrus.Info("Store cache disabled")
		cache = storecache.NewNullCache()
	}

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			store := mongo.NewStore(client.Database("main"), cache)
			ctx := apicontext.NewContext(store, c)

			return next(ctx)
		}
	})

	// Public routes for external access through API gateway
	publicAPI := e.Group("/api")

	// Internal routes only accessible by other services in the local container network
	internalAPI := e.Group("/internal")

	internalAPI.GET(routes.AuthRequestURL, apicontext.Handler(routes.AuthRequest), apicontext.Middleware(routes.AuthMiddleware))
	publicAPI.POST(routes.AuthDeviceURL, apicontext.Handler(routes.AuthDevice))
	publicAPI.POST(routes.AuthDeviceURLV2, apicontext.Handler(routes.AuthDevice))
	publicAPI.POST(routes.AuthUserURL, apicontext.Handler(routes.AuthUser))
	publicAPI.POST(routes.AuthUserURLV2, apicontext.Handler(routes.AuthUser))
	publicAPI.GET(routes.AuthUserURLV2, apicontext.Handler(routes.AuthUserInfo))
	internalAPI.GET(routes.AuthUserTokenURL, apicontext.Handler(routes.AuthGetToken))
	publicAPI.POST(routes.AuthPublicKeyURL, apicontext.Handler(routes.AuthPublicKey))
	publicAPI.GET(routes.AuthUserTokenURL, apicontext.Handler(routes.AuthSwapToken))

	publicAPI.PATCH(routes.UpdateUserDataURL, apicontext.Handler(routes.UpdateUserData))
	publicAPI.PATCH(routes.UpdateUserPasswordURL, apicontext.Handler(routes.UpdateUserPassword))
	publicAPI.PUT(routes.EditSessionRecordStatusURL, apicontext.Handler(routes.EditSessionRecordStatus))
	publicAPI.GET(routes.GetSessionRecordURL, apicontext.Handler(routes.GetSessionRecord))

	publicAPI.GET(routes.GetDeviceListURL,
		middlewares.Authorize(apicontext.Handler(routes.GetDeviceList)))
	publicAPI.GET(routes.GetDeviceURL,
		middlewares.Authorize(apicontext.Handler(routes.GetDevice)))
	publicAPI.DELETE(routes.DeleteDeviceURL, apicontext.Handler(routes.DeleteDevice))
	publicAPI.PATCH(routes.RenameDeviceURL, apicontext.Handler(routes.RenameDevice))
	internalAPI.POST(routes.OfflineDeviceURL, apicontext.Handler(routes.OfflineDevice))
	internalAPI.GET(routes.LookupDeviceURL, apicontext.Handler(routes.LookupDevice))
	publicAPI.PATCH(routes.UpdateStatusURL, apicontext.Handler(routes.UpdatePendingStatus))
	publicAPI.GET(routes.GetSessionsURL,
		middlewares.Authorize(apicontext.Handler(routes.GetSessionList)))
	publicAPI.GET(routes.GetSessionURL,
		middlewares.Authorize(apicontext.Handler(routes.GetSession)))
	internalAPI.PATCH(routes.SetSessionAuthenticatedURL, apicontext.Handler(routes.SetSessionAuthenticated))
	internalAPI.POST(routes.CreateSessionURL, apicontext.Handler(routes.CreateSession))
	internalAPI.POST(routes.FinishSessionURL, apicontext.Handler(routes.FinishSession))
	internalAPI.POST(routes.RecordSessionURL, apicontext.Handler(routes.RecordSession))
	publicAPI.GET(routes.PlaySessionURL, apicontext.Handler(routes.PlaySession))
	publicAPI.DELETE(routes.RecordSessionURL, apicontext.Handler(routes.DeleteRecordedSession))

	publicAPI.GET(routes.GetStatsURL,
		middlewares.Authorize(apicontext.Handler(routes.GetStats)))

	publicAPI.GET(routes.GetPublicKeysURL, apicontext.Handler(routes.GetPublicKeys))
	publicAPI.POST(routes.CreatePublicKeyURL, apicontext.Handler(routes.CreatePublicKey))
	publicAPI.PUT(routes.UpdatePublicKeyURL, apicontext.Handler(routes.UpdatePublicKey))
	publicAPI.DELETE(routes.DeletePublicKeyURL, apicontext.Handler(routes.DeletePublicKey))
	internalAPI.GET(routes.GetPublicKeyURL, apicontext.Handler(routes.GetPublicKey))
	internalAPI.POST(routes.CreatePrivateKeyURL, apicontext.Handler(routes.CreatePrivateKey))
	internalAPI.POST(routes.EvaluateKeyURL, apicontext.Handler(routes.EvaluateKeyHostname))

	publicAPI.GET(routes.ListNamespaceURL, apicontext.Handler(routes.GetNamespaceList))
	publicAPI.GET(routes.GetNamespaceURL, apicontext.Handler(routes.GetNamespace))
	publicAPI.POST(routes.CreateNamespaceURL, apicontext.Handler(routes.CreateNamespace))
	publicAPI.DELETE(routes.DeleteNamespaceURL, apicontext.Handler(routes.DeleteNamespace))
	publicAPI.PUT(routes.EditNamespaceURL, apicontext.Handler(routes.EditNamespace))
	publicAPI.PATCH(routes.AddNamespaceUserURL, apicontext.Handler(routes.AddNamespaceUser))
	publicAPI.PATCH(routes.RemoveNamespaceUserURL, apicontext.Handler(routes.RemoveNamespaceUser))

	publicAPI.GET(routes.ListTokenURL, apicontext.Handler(routes.ListToken))
	publicAPI.GET(routes.GetTokenURL, apicontext.Handler(routes.GetToken))
	publicAPI.POST(routes.CreateTokenURL, apicontext.Handler(routes.CreateToken))
	publicAPI.DELETE(routes.DeleteTokenURL, apicontext.Handler(routes.DeleteToken))
	publicAPI.PATCH(routes.UpdateTokenURL, apicontext.Handler(routes.UpdateToken))

	e.Logger.Fatal(e.Start(":8080"))

	return nil
}
