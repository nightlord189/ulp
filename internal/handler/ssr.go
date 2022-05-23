package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/nightlord189/ulp/internal/model"
	"net/http"
	"time"
)

func (h *Handler) Index(c echo.Context) error {
	role := model.Role(c.Get("role").(string))
	switch role {
	case model.RoleStudent:
		return c.Redirect(302, "/attempts")
	case model.RoleTutor:
		return c.Redirect(302, "/tasks")
	}
	return c.Render(http.StatusBadRequest, "message.html", model.NewTmplMessage("Ошибка определения роли", true))
}

func (h *Handler) GetTasks(c echo.Context) error {
	role := model.Role(c.Get("role").(string))
	if role != model.RoleTutor {
		return c.Redirect(302, "/")
	}
	userID := fmt.Sprintf("%v", c.Get("user_id"))
	resp, err := h.Service.GetTasks(userID)
	if err != nil {
		return c.Render(http.StatusBadRequest, "message.html", model.NewTmplMessage("Ошибка загрузки заданий: "+err.Error(), true))
	}
	return c.Render(http.StatusOK, "tasks.html", resp)
}

func (h *Handler) GetAttempts(c echo.Context) error {
	role := model.Role(c.Get("role").(string))
	if role != model.RoleStudent {
		return c.Redirect(302, "/")
	}
	userID := fmt.Sprintf("%v", c.Get("user_id"))
	resp, err := h.Service.GetAttempts(userID)
	if err != nil {
		return c.Render(http.StatusBadRequest, "message.html", model.NewTmplMessage("Ошибка загрузки решений: "+err.Error(), true))
	}
	return c.Render(http.StatusOK, "attempts.html", resp)
}

func (h *Handler) GetAttempt(c echo.Context) error {
	authorized := getBool(c, "authorized")
	var role model.Role
	if authorized {
		role = model.Role(c.Get("role").(string))
	}
	id := c.Param("id")
	resp, err := h.Service.GetAttempt(id, authorized, string(role))
	if err != nil {
		return c.Render(http.StatusBadRequest, "message.html", model.NewTmplMessage("Ошибка получения решения: "+err.Error(), true))
	}
	return c.Render(http.StatusOK, "attempt.html", resp)
}

func (h *Handler) Logout(c echo.Context) error {
	removeCookie(c, "auth")
	return c.Redirect(302, "/")
}

func (h *Handler) GetAuth(c echo.Context) error {
	authorized := getBool(c, "authorized")
	if !authorized {
		return c.Render(http.StatusOK, "auth.html", "")
	}
	return c.Redirect(302, "/")
}

func (h *Handler) GetReg(c echo.Context) error {
	authorized := getBool(c, "authorized")
	if !authorized {
		return c.Render(http.StatusOK, "reg.html", "")
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
