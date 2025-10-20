package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rod1kutzyy/OnTrack/internal/delivery/http/handlers"
)

func RegisterTodoRoutes(r *gin.Engine, h *handlers.TodoHandler) {
	todos := r.Group("/todos")
	{
		todos.POST("/", h.Create)
		todos.GET("/", h.GetAll)
		todos.GET("/:id", h.GetByID)
		todos.PUT("/:id", h.Update)
		todos.DELETE("/:id", h.Delete)
	}
}
