package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nightlord189/ulp/internal/config"
	"github.com/nightlord189/ulp/internal/db"
	"github.com/nightlord189/ulp/internal/handler"
	"github.com/nightlord189/ulp/internal/service"
	"net/http"
)

func main() {
	fmt.Println("start")

	cfg := config.Load("configs/config.json")
	dbInstance, err := db.InitDb(cfg)
	if err != nil {
		panic(fmt.Sprintf("error init db: %v", err))
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

	e.GET("/", hlr.Index, hlr.CheckCookieJwtMiddleware)
	e.GET("/auth", hlr.GetAuth, hlr.CheckCookieJwtMiddleware)
	e.GET("/reg", hlr.GetReg, hlr.CheckCookieJwtMiddleware)
	e.POST("/auth", hlr.PostAuth)
	e.POST("/reg", hlr.PostReg)

	err = e.Start(fmt.Sprintf(":%d", cfg.HttpPort))
	if err != nil {
		panic(fmt.Sprintf("error start server: %v", err))
	}
}
