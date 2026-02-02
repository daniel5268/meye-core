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
	UserHandler *handler.UserHandler
	AuthHandler *handler.AuthHandler
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
}

func (r *Router) setupUserRoutes(group *gin.RouterGroup) {
	users := group.Group("/users")
	{
		users.POST("", r.handlers.AuthHandler.InternalAPIKeyMiddleware(), r.handlers.UserHandler.CreateUser)
		users.POST("/login", r.handlers.UserHandler.Login)
	}
}

// Engine returns the Gin engine
func (r *Router) Engine() *gin.Engine {
	return r.engine
}
