package main

import (
	"net/http"

	"github.com/motoki317/go-with-world-database/database"
	"github.com/motoki317/go-with-world-database/login"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := database.SetUpDatabase()
	store := database.SetUpSessionDatabase(db)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	login.SetUpLoginRoutes(e, db)

	withLogin := e.Group("")
	withLogin.Use(login.CheckLogin)
	withLogin.GET("/cities/:cityName", database.MakeRetrieveCityHandler(db))
	withLogin.GET("/whoami", login.WhoAmI)

	e.Start(":10901")
}
