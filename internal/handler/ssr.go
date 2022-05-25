package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/nightlord189/ulp/internal/model"
	"net/http"
	"strconv"
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
	return renderMessage(c, "Ошибка определения роли", true)
}

func (h *Handler) GetTasks(c echo.Context) error {
	userID := getString(c, "user_id")
	resp, err := h.Service.GetTasks(userID)
	if err != nil {
		return renderMessage(c, "Ошибка загрузки заданий: "+err.Error(), true)
	}
	return c.Render(http.StatusOK, "tasks.html", resp)
}

func (h *Handler) GetCreateTask(c echo.Context) error {
	data, _ := h.Service.GetCreateTask()
	return c.Render(http.StatusOK, "edit_task.html", data)
}

func (h *Handler) GetAttemptTask(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return renderMessage(c, "Некорректный id задания: "+err.Error(), true)
	}
	data, _ := h.Service.GetAttemptTask(id)
	return c.Render(http.StatusOK, "upload_attempt.html", data)
}

func (h *Handler) GetEditTask(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return renderMessage(c, "Некорректный id задания: "+err.Error(), true)
	}
	response, err := h.Service.GetEditTask(id)
	if err != nil {
		return renderMessage(c, "Ошибка: "+err.Error(), true)
	}
	return c.Render(http.StatusOK, "edit_task.html", response)
}

func (h *Handler) PostCreateTask(c echo.Context) error {
	var req model.ChangeTaskRequest
	err := c.Bind(&req)
	if err != nil {
		return renderMessage(c, "Ошибка формы: "+err.Error(), true)
	}
	if err = req.IsValid(); err != nil {
		return renderMessage(c, "Некорректный запрос: "+err.Error(), true)
	}
	userIDStr := getString(c, "user_id")
	userID, _ := strconv.Atoi(userIDStr)
	req.CreatorID = userID
	err = h.Service.CreateTask(req)
	if err != nil {
		return renderMessage(c, "Не удалось создать задание: "+err.Error(), true)
	}
	return renderMessage(c, "Задание успешно создано", false)
}

func (h *Handler) PostEditTask(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return renderMessage(c, "Некорректный id задания: "+err.Error(), true)
	}
	var req model.ChangeTaskRequest
	err = c.Bind(&req)
	if err != nil {
		return renderMessage(c, "Ошибка формы: "+err.Error(), true)
	}
	if err = req.IsValid(); err != nil {
		return renderMessage(c, "Некорректный запрос: "+err.Error(), true)
	}
	req.ID = id
	userIDStr := getString(c, "user_id")
	userID, _ := strconv.Atoi(userIDStr)
	req.CreatorID = userID
	err = h.Service.EditTask(req)
	if err != nil {
		return renderMessage(c, "Не удалось изменить задание: "+err.Error(), true)
	}
	return renderMessage(c, "Задание успешно обновлено", false)
}

func (h *Handler) DeleteTask(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return renderMessage(c, "Некорректный id задания: "+err.Error(), true)
	}
	err = h.Service.DeleteTask(id)
	if err != nil {
		return renderMessage(c, "Ошибка удаления задания: "+err.Error(), true)
	}
	return c.Redirect(302, "/tasks")
}

func (h *Handler) GetAttempts(c echo.Context) error {
	userID := fmt.Sprintf("%v", c.Get("user_id"))
	resp, err := h.Service.GetAttempts(userID)
	if err != nil {
		return renderMessage(c, "Ошибка загрузки решений: "+err.Error(), true)
	}
	return c.Render(http.StatusOK, "attempts.html", resp)
}

func (h *Handler) GetAttempt(c echo.Context) error {
	authorized := getBool(c, "authorized")
	var role model.Role
	if authorized {
		role = model.Role(getString(c, "role"))
	}
	id := c.Param("id")
	resp, err := h.Service.GetAttempt(id, authorized, string(role))
	if err != nil {
		return renderMessage(c, "Ошибка получения решения: "+err.Error(), true)
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
		return renderMessage(c, "Ошибка формы: "+err.Error(), true)
	}
	if err = req.IsValid(); err != nil {
		return renderMessage(c, "Некорректный запрос: "+err.Error(), true)
	}
	token, err := h.Service.Auth(req)
	if err != nil {
		return renderMessage(c, "Ошибка авторизации: "+err.Error(), true)
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
		return renderMessage(c, "Ошибка формы: "+err.Error(), true)
	}
	if err = req.IsValid(); err != nil {
		return renderMessage(c, "Некорректный запрос: "+err.Error(), true)
	}
	err = h.Service.Reg(req)
	if err != nil {
		return renderMessage(c, "Ошибка регистрации: "+err.Error(), true)
	}
	return renderMessage(c, "Регистрация прошла успешно", false)
}
