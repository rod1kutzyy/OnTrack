package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rod1kutzyy/OnTrack/internal/config"
	"github.com/rod1kutzyy/OnTrack/internal/handler"
	"github.com/rod1kutzyy/OnTrack/internal/middleware"
)

func SetupRouter(cfg *config.Config, todoHandler *handler.TodoHandler) *gin.Engine {
	if cfg.Logger.Level == "debug" || cfg.Logger.Level == "trace" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})

	router.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
			"time":   time.Now().UTC(),
		})
	})

	v1 := router.Group("/api/v1")
	{
		todos := v1.Group("/todos")
		{
			todos.POST("", todoHandler.CreateTodo)
			todos.GET("", todoHandler.GetAllTodos)
			todos.GET("/:id", todoHandler.GetTodoByID)
			todos.PUT("/:id", todoHandler.UpdateTodo)
			todos.DELETE("/:id", todoHandler.DeleteTodo)
			todos.PATCH("/:id/toggle", todoHandler.ToggleTodoComplete)
		}
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Not found",
			"message": "The requested endpoint does not exist",
			"path":    c.Request.URL.Path,
		})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"success": false,
			"error":   "Method Not Allowed",
			"message": "The HTTP method is not supported for this endpoint",
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
		})
	})

	return router
}
