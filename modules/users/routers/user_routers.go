package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/yuta_2710/go-clean-arc-reviews/middleware"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/users/handlers"
)

func InitUserRouters(http handlers.UserHandler, mdwr echo.MiddlewareFunc, routes *echo.Group) {
	admin := routes.Group("/users")
	admin.Use(mdwr)
	admin.Use(middleware.IsAdmin())

	admin.GET("/", http.GetUsers)
	admin.GET("/:id", http.GetUserById)
	admin.POST("/", http.CreateNewUser)
}
