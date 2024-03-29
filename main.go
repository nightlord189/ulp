package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nightlord189/ulp/internal/config"
	"github.com/nightlord189/ulp/internal/db"
	"github.com/nightlord189/ulp/internal/handler"
	"github.com/nightlord189/ulp/internal/service"
	"github.com/nightlord189/ulp/pkg/log"
	"github.com/rs/zerolog"
	stdLog "log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("start #1")

	cfg := config.Load("configs/config.json")

	if err := log.InitLogger(cfg.LogLevel, "ulp", ""); err != nil {
		stdLog.Fatalf("error on init logger: %v", err)
	}

	ctx := context.Background()

	zerolog.Ctx(ctx).Info().Msg("start #2")

	if _, err := os.Stat(cfg.AttemptsPath); os.IsNotExist(err) {
		if err := os.Mkdir(cfg.AttemptsPath, os.ModePerm); err != nil {
			zerolog.Ctx(ctx).Fatal().Msgf("error create attempts directory: %v", err)
		}
	}

	dbInstance, err := db.InitDb(cfg)
	if err != nil {
		zerolog.Ctx(ctx).Fatal().Msgf("error init db: %v", err)
	}

	srv := service.NewService(cfg, dbInstance)
	hlr := handler.NewHandler(cfg, dbInstance, srv, cfg.TemplatesPath)

	e := echo.New()
	e.Debug = cfg.HttpDebug
	e.Renderer = hlr

	e.Static("/static", "web/static")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "healthz")
	})

	e.GET("/", hlr.Index, hlr.CookieJwtMiddleware, hlr.AuthorizedMiddleware)
	e.GET("/logout", hlr.Logout)
	e.GET("/auth", hlr.GetAuth, hlr.CookieJwtMiddleware)
	e.GET("/reg", hlr.GetReg, hlr.CookieJwtMiddleware)

	e.GET("/tasks", hlr.GetTasks, hlr.CookieJwtMiddleware, hlr.AuthorizedMiddleware, hlr.TutorMiddleware)
	e.GET("/task/create", hlr.GetCreateTask, hlr.CookieJwtMiddleware, hlr.AuthorizedMiddleware, hlr.TutorMiddleware)
	e.POST("/task/create", hlr.PostCreateTask, hlr.CookieJwtMiddleware, hlr.AuthorizedMiddleware, hlr.TutorMiddleware)
	e.GET("/task/:id/edit", hlr.GetEditTask, hlr.CookieJwtMiddleware, hlr.AuthorizedMiddleware, hlr.TutorMiddleware)
	e.POST("/task/:id/edit", hlr.PostEditTask, hlr.CookieJwtMiddleware, hlr.AuthorizedMiddleware, hlr.TutorMiddleware)
	e.POST("/task/:id/delete", hlr.DeleteTask, hlr.CookieJwtMiddleware, hlr.AuthorizedMiddleware, hlr.TutorMiddleware)
	e.GET("/task/:id/attempt", hlr.GetAttemptTask, hlr.CookieJwtMiddleware, hlr.AuthorizedMiddleware, hlr.StudentMiddleware)
	e.POST("/task/:id/attempt", hlr.PostAttempt, hlr.CookieJwtMiddleware, hlr.AuthorizedMiddleware, hlr.StudentMiddleware)

	e.GET("/attempts", hlr.GetAttempts, hlr.CookieJwtMiddleware, hlr.AuthorizedMiddleware, hlr.StudentMiddleware)
	e.GET("/attempt/:id", hlr.GetAttempt, hlr.CookieJwtMiddleware)

	e.POST("/auth", hlr.PostAuth)
	e.POST("/reg", hlr.PostReg)

	e.GET("/api/template/dockerfile", hlr.GetDockerfileTemplates)

	fmt.Println("server arch:", cfg.ContainerArchitecture)

	zerolog.Ctx(ctx).Info().Msg("start #3 (starting web-server)")

	err = e.Start(fmt.Sprintf(":%d", cfg.HttpPort))
	if err != nil {
		zerolog.Ctx(ctx).Fatal().Msgf("error start server: %v", err)
	}
}
