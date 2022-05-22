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
		_, err = service.ValidateJwtToken(authCookie.Value, h.Config.Auth.Secret)
		if err != nil {
			// remove cookie
			/*c.Set("authorized", false)
			cookie := new(http.Cookie)
			cookie.Name = "auth"
			cookie.Value = ""
			cookie.Expires = time.Now()
			c.SetCookie(cookie)*/
			fmt.Println("invalid token middleware", err, authCookie.Value)
			return next(c)
		}
		c.Set("authorized", true)
		return next(c)
	}
}
