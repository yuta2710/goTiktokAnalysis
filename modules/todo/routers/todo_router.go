package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/todo/handlers"
)

func InitTodoRoutes(http handlers.TodoHandler, mdwr echo.MiddlewareFunc, routes *echo.Group) {
	protected := routes.Group("/todos")
	protected.Use(mdwr)

	protected.GET("/:authId", http.FindAllByUserId)
	protected.GET("/:id", http.FindById)
	protected.POST("/", http.CreateNewTodo)
}
