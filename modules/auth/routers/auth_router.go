package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/auth/handlers"
)

func InitAuthRouters(http handlers.AuthHandler, protectMdwr echo.MiddlewareFunc, routes *echo.Group) {
	// ROUTE_VAL := "/auth"
	publicRoutes := routes.Group("/auth")
	publicRoutes.POST("/login", http.Login)
	publicRoutes.POST("/signup", http.SignUp)

	protectedRoutes := routes.Group("/auth")
	protectedRoutes.Use(protectMdwr)

	protectedRoutes.GET("/profile", http.Profile)
}
