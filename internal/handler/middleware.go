package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/nightlord189/ulp/internal/service"
)

func (h *Handler) CheckCookieJwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("middleware")
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
		c.Set("user_id", claims["user_id"])
		c.Set("username", claims["username"])
		c.Set("role", claims["role"])
		return next(c)
	}
}

func (h *Handler) RedirectUnauthorizedMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("middleware2")
		authorized := getBool(c, "authorized")
		if !authorized {
			return c.Redirect(302, "/auth")
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
