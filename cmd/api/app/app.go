package app

import (
	"fmt"
	"net/http"
	"time"

	logger "github.com/Abhishekjha321/community_service/log"
	"github.com/gin-contrib/cors"

	"github.com/Abhishekjha321/community_service/storage/cache" // write this code locally
	api "github.com/Abhishekjha321/community_service/api/http/v1"
	"github.com/Abhishekjha321/community_service/internal/logic/community/model"
	"github.com/Abhishekjha321/community_service/internal/logic/community/repo"
	"github.com/Abhishekjha321/community_service/internal/logic/community/service"
	"github.com/Abhishekjha321/community_service/pkg/config"
	"github.com/Abhishekjha321/community_service/pkg/store/db"
	"github.com/Abhishekjha321/community_service/pkg/store/redis"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type services struct {
	// application services
	communityService model.Service
}

type controller struct {
	// app controllers
	communityController api.CommunityController
}

type Application struct {
	db         *db.Store
	cache      cache.CacheBase
	services   services
	controller controller
	router     *gin.Engine
	http       *http.Server
}

func (a *Application) initStores() {
	var err error
	a.db, err = db.NewPostgresStorage()
	if err != nil {
		panic(fmt.Errorf("db initialization failed: %w", err))
	}
}

func (a *Application) initCache() {
	var err error
	a.cache, err = redis.InitializeRedis()
	if err != nil {
		panic(err)
	}
}

func (a *Application) initServices() {
	a.services.communityService = service.NewService(repo.NewRepo(a.db), a.cache)
}

func (a *Application) initControllers() {
	a.controller.communityController = api.NewCommunityController(a.services.communityService)
}

func (a *Application) setUpHandlers() *gin.Engine {
	router := gin.Default()

	// Apply CORS middleware globally
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type","x-user-id"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Set up OpenTelemetry middleware
	router.Use(otelgin.Middleware(config.Config.Name))

	// Add public and private routes
	api.AddPublicRoutes(router, a.controller.communityController)
	api.AddPrivateRoutes(router, a.controller.communityController)

	return router
}

func (a *Application) Init() {
	a.initStores()
	a.initCache()
	a.initServices()
	a.initControllers()
	a.router = a.setUpHandlers()
	a.http = &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Config.Server.Port),
		Handler:      a.router,
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
		IdleTimeout:  time.Second * 60,
	}
}

func (a *Application) Start() {
	defer logger.GetLogger().Errorf("stopped http server")
	fmt.Printf("server is listening on port: %d \n", config.Config.Server.Port)
	if err := a.http.ListenAndServe(); err != nil {
		logger.GetLogger().WithError(err).Fatal("failed to start http server")
	}
}
