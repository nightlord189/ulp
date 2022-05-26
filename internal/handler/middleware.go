package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/nightlord189/ulp/internal/model"
	"github.com/nightlord189/ulp/internal/service"
	"strconv"
)

func (h *Handler) CookieJwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authCookie, err := c.Cookie("auth")
		if err != nil || authCookie.Value == "" {
			c.Set("authorized", false)
			fmt.Println("no cookie")
			return next(c)
		}
		claims, err := service.ValidateJwtToken(authCookie.Value, h.Config.Auth.Secret)
		if err != nil {
			// remove cookie
			c.Set("authorized", false)
			removeCookie(c, "auth")
			fmt.Println("invalid token middleware", err, authCookie.Value)
			return next(c)
		}
		c.Set("authorized", true)
		// проставляем нужные данные для дальнейшего использования в обработчиках
		userIDStr := fmt.Sprintf("%v", claims["user_id"])
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return fmt.Errorf("user_id has bad format: %w", err)
		}
		c.Set("user_id", userID)
		c.Set("username", claims["username"])
		c.Set("role", claims["role"])
		return next(c)
	}
}

func (h *Handler) AuthorizedMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorized := getBool(c, "authorized")
		if !authorized {
			return c.Redirect(302, "/auth")
		}
		return next(c)
	}
}

func (h *Handler) TutorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := model.Role(c.Get("role").(string))
		if role != model.RoleTutor {
			return renderMessage(c, "Данная страница предназначена только для преподавателей", true)
		}
		return next(c)
	}
}

func (h *Handler) StudentMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := model.Role(c.Get("role").(string))
		if role != model.RoleStudent {
			return renderMessage(c, "Данная страница предназначена только для студентов", true)
		}
		return next(c)
	}
}

func (h *Handler) Middleware3(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("middleware3")
		authorized := getBool(c, "authorized")
		fmt.Println("authorized", authorized)
		return next(c)
	}
}
