package api

import (
	"meye-core/internal/infrastructure/api/handler"
	customValidator "meye-core/internal/infrastructure/api/validator"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Router handles HTTP routing and middleware configuration
type Router struct {
	engine   *gin.Engine
	handlers *Handlers
}

// RouterConfig holds dependencies needed for routing
type Handlers struct {
	UserHandler     *handler.UserHandler
	AuthHandler     *handler.AuthHandler
	CampaignHandler *handler.CampaignHandler
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		customValidator.RegisterCustomValidators(v)
	}
}

func NewRouter(handlers *Handlers) *Router {
	engine := gin.Default()
	engine.Use(gin.Recovery())

	router := &Router{
		engine:   engine,
		handlers: handlers,
	}

	router.setupRoutes()

	return router
}

func (r *Router) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"service": "meye-core",
	})
}

func (r *Router) setupRoutes() {
	r.engine.GET("/health", r.healthCheck)

	v1 := r.engine.Group("/api/v1")
	r.setupUserRoutes(v1)
	r.setupCampaignRoutes(v1)
	r.setupPjRoutes(v1)
}

func (r *Router) setupUserRoutes(group *gin.RouterGroup) {
	users := group.Group("/users")
	{
		users.POST("",
			r.handlers.AuthHandler.AuthMiddleware(),
			r.handlers.AuthHandler.RequireAdminRole(),
			r.handlers.UserHandler.CreateUser,
		)
		users.POST("/login", r.handlers.UserHandler.Login)
		users.GET("/players",
			r.handlers.AuthHandler.AuthMiddleware(),
			r.handlers.AuthHandler.RequireMasterRole(),
			r.handlers.UserHandler.GetPlayers,
		)
	}
}

func (r *Router) setupCampaignRoutes(group *gin.RouterGroup) {
	campaigns := group.Group("/campaigns")
	campaigns.Use(r.handlers.AuthHandler.AuthMiddleware())
	{
		campaigns.POST("",
			r.handlers.AuthHandler.RequireMasterRole(),
			r.handlers.CampaignHandler.CreateCampaign,
		)
		campaigns.POST("/:campaignID/invitations",
			r.handlers.AuthHandler.RequireCampaignMaster(),
			r.handlers.CampaignHandler.InviteUser,
		)
		campaigns.POST("/:campaignID/pjs",
			r.handlers.AuthHandler.RequirePlayerRole(),
			r.handlers.CampaignHandler.CreatePJ,
		)
		campaigns.POST("/:campaignID/sessions",
			r.handlers.AuthHandler.RequireCampaignMaster(),
			r.handlers.CampaignHandler.CreateSession,
		)
		campaigns.GET("/:campaignID",
			r.handlers.AuthHandler.RequireCampaignMaster(),
			r.handlers.CampaignHandler.GetCampaign,
		)
	}
}

func (r *Router) setupPjRoutes(group *gin.RouterGroup) {
	pjs := group.Group("/pjs")
	pjs.Use(r.handlers.AuthHandler.AuthMiddleware())
	pjs.Use(r.handlers.AuthHandler.RequirePjUser())
	{
		pjs.PUT("/:pjID/stats", r.handlers.CampaignHandler.UpdatePJStats)
		pjs.GET("/:pjID", r.handlers.CampaignHandler.GetPj)
	}
}

// Engine returns the Gin engine
func (r *Router) Engine() *gin.Engine {
	return r.engine
}
