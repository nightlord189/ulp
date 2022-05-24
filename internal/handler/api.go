package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) GetDockerfileTemplates(c echo.Context) error {
	entities, err := h.Service.GetDockerfileTemplates()
	if err != nil {
		return c.String(http.StatusBadRequest, "Error get templates: "+err.Error())
	}
	return c.JSON(http.StatusOK, entities)
}
