package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/nightlord189/ulp/internal/config"
	"github.com/nightlord189/ulp/internal/db"
	"github.com/nightlord189/ulp/internal/service"
	"html/template"
	"io"
	"net/http"
	"time"
)

// Handler - структура со ссылками на зависимости
type Handler struct {
	Config    *config.Config
	DB        *db.Manager
	Service   *service.Service
	templates *template.Template
}

// NewHandler - конструктор Handler
func NewHandler(
	cfg *config.Config,
	db *db.Manager,
	srv *service.Service,
	pathToTemplates string,
) *Handler {
	templates := template.Must(template.ParseGlob(pathToTemplates))
	handler := Handler{
		Config:    cfg,
		DB:        db,
		Service:   srv,
		templates: templates,
	}
	return &handler
}

func (h *Handler) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return h.templates.ExecuteTemplate(w, name, data)
}

func getBool(c echo.Context, key string) bool {
	var value bool
	valueCommon := c.Get(key)
	if valueCommon != nil {
		value = valueCommon.(bool)
	}
	return value
}

func removeCookie(c echo.Context, name string) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = ""
	cookie.Expires = time.Now()
	c.SetCookie(cookie)
}
