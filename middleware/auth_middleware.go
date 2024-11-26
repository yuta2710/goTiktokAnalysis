package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/users/entities"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/users/repositories"
	"github.com/yuta_2710/go-clean-arc-reviews/shared"
)

// Protect middleware for authorization
func Protect(getUserById func(authId string) (*entities.FetchUserDto, error)) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" {
				return shared.Response(c, false, http.StatusUnauthorized, "Authorization token is missing", nil, nil)
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return shared.Response(c, false, http.StatusUnauthorized, "Invalid authorization format", nil, nil)
			}

			acToken := parts[1]
			acSecret := os.Getenv("ACCESS_SECRET")

			// Parse and validate token
			token, err := jwt.Parse(acToken, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(acSecret), nil
			})

			if err != nil {
				return shared.Response(c, false, http.StatusUnauthorized, "Unauthorized", nil, nil)
			}

			// Extract claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				return shared.Response(c, false, http.StatusUnauthorized, "Invalid token claims", nil, nil)
			}

			authId, ok := claims["authId"].(string)
			if !ok {
				return shared.Response(c, false, http.StatusUnauthorized, "User ID not found in token", nil, nil)
			}

			fmt.Printf("Auth id %s", authId)

			// Use the provided function to fetch the user
			user, err := getUserById(authId)
			if err != nil {
				return shared.Response(c, false, http.StatusUnauthorized, "User not found", nil, nil)
			}

			// Attach user info to context
			c.Set("user", user)

			ctx := context.WithValue(c.Request().Context(), "authId", authId)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func IsAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("[DEBUG] IsAdmin middleware triggered")
			user, ok := c.Get("user").(*entities.FetchUserDto)
			fmt.Println(user)
			fmt.Printf("Role cua thang nay la %s", user.Role)

			if !ok {
				return shared.Response(c, false, http.StatusUnauthorized, "[UNAUTHORIZED]: User ID not found", nil, nil)
			}

			if user.Role != "admin" {
				return shared.Response(c, false, http.StatusForbidden, "Forbidden: Admin role required", nil, nil)
			}

			return next(c)
		}
	}
}

func NewProtectMiddleware(userRepo repositories.UserRepository) echo.MiddlewareFunc {
	return Protect(func(authId string) (*entities.FetchUserDto, error) {
		// Decode the authId
		decodedUID, err := shared.DecomposeUidV2(authId)
		if err != nil {
			return nil, fmt.Errorf("invalid authId: %v", err)
		}

		idStr := fmt.Sprintf("%d", decodedUID.GetLocalID())
		id, _ := strconv.Atoi(idStr)
		// Find user by ID
		user, err := userRepo.FindById(id)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Role cua thang nay trong protect middelware la %s", user.Role)

		fmt.Println("Hahaha")

		// Map user to FetchUserDto
		return &entities.FetchUserDto{
			FakeId:    authId,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      user.Role,
		}, nil
	})
}
