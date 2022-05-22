package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/nightlord189/ulp/internal/model"
	"net/http"
	"time"
)

func (h *Handler) Index(c echo.Context) error {
	authorized := h.GetBool(c, "authorized")
	if authorized {
		return c.String(http.StatusOK, "authorized3")
	}
	return c.Redirect(302, "/auth")
}

func (h *Handler) GetAuth(c echo.Context) error {
	authorized := h.GetBool(c, "authorized")
	if !authorized {
		return c.Render(http.StatusOK, "auth.html", "World")
	}
	return c.Redirect(302, "/")
}

func (h *Handler) PostAuth(c echo.Context) error {
	var req model.AuthRequest
	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusUnprocessableEntity, "error bind request: "+err.Error())
	}
	if err = req.IsValid(); err != nil {
		return c.String(http.StatusUnprocessableEntity, "invalid request: "+err.Error())
	}
	token, err := h.Service.Auth(req)
	if err != nil {
		return c.String(http.StatusBadRequest, "err auth: "+err.Error())
	}
	cookie := new(http.Cookie)
	cookie.Name = "auth"
	cookie.Value = token
	cookie.Expires = time.Now().Add(time.Second * time.Duration(h.Config.Auth.ExpTime))
	c.SetCookie(cookie)
	fmt.Println("token issued", token)
	return c.Redirect(302, "/")
}
