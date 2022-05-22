package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/nightlord189/ulp/internal/model"
	"net/http"
	"time"
)

func (h *Handler) Index(c echo.Context) error {
	authorized := getBool(c, "authorized")
	if authorized {
		username := c.Get("username")
		return c.String(http.StatusOK, "authorized4 "+fmt.Sprintf("%v", username))
	}
	return c.Redirect(302, "/auth")
}

func (h *Handler) GetAuth(c echo.Context) error {
	authorized := getBool(c, "authorized")
	if !authorized {
		return c.Render(http.StatusOK, "auth.html", "World")
	}
	return c.Redirect(302, "/")
}

func (h *Handler) GetReg(c echo.Context) error {
	authorized := getBool(c, "authorized")
	if !authorized {
		return c.Render(http.StatusOK, "reg.html", "World")
	}
	return c.Redirect(302, "/")
}

func (h *Handler) PostAuth(c echo.Context) error {
	var req model.AuthRequest
	err := c.Bind(&req)
	if err != nil {
		return c.Render(http.StatusUnprocessableEntity, "message.html", model.NewTmplMessage("Ошибка формы: "+err.Error(), true))
	}
	if err = req.IsValid(); err != nil {
		return c.Render(http.StatusUnprocessableEntity, "message.html", model.NewTmplMessage("Некорректный запрос: "+err.Error(), true))
	}
	token, err := h.Service.Auth(req)
	if err != nil {
		return c.Render(http.StatusBadRequest, "message.html", model.NewTmplMessage("Ошибка авторизации: "+err.Error(), true))
	}
	cookie := new(http.Cookie)
	cookie.Name = "auth"
	cookie.Value = token
	cookie.Expires = time.Now().Add(time.Second * time.Duration(h.Config.Auth.ExpTime))
	c.SetCookie(cookie)
	fmt.Println("token issued", token)
	return c.Redirect(302, "/")
}

func (h *Handler) PostReg(c echo.Context) error {
	var req model.RegRequest
	err := c.Bind(&req)
	if err != nil {
		return c.Render(http.StatusUnprocessableEntity, "message.html", model.NewTmplMessage("Ошибка формы: "+err.Error(), true))
	}
	if err = req.IsValid(); err != nil {
		return c.Render(http.StatusUnprocessableEntity, "message.html", model.NewTmplMessage("Некорректный запрос: "+err.Error(), true))
	}
	err = h.Service.Reg(req)
	if err != nil {
		return c.Render(http.StatusBadRequest, "message.html", model.NewTmplMessage("Ошибка регистрации: "+err.Error(), true))
	}
	return c.Render(http.StatusOK, "message.html", model.NewTmplMessage("Регистрация прошла успешно. ", false))
}
