package database

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

type City struct {
	ID          int    `json:"id,omitempty" db:"ID"`
	Name        string `json:"name,omitempty"  db:"Name"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

type Country struct {
	Name       string `json:"name,omitempty" db:"Name"`
	Population int    `json:"population,omitempty" db:"Population"`
}

type payloadCity struct {
	payloadUser
	City City `json:"city,omitempty"`
}

type payloadUser struct {
	Username string `json:"username,omitempty"`
}

func createPayloadUser(username string) payloadUser {
	return payloadUser{
		Username: username,
	}
}

// MakeRetrieveCityHandler Worldデータベースから都市情報を取得します
func MakeRetrieveCityHandler(db *sqlx.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		cityName := c.Param("cityName")

		city := City{}
		db.Get(&city, "SELECT * FROM city WHERE Name=?", cityName)
		if city.Name == "" {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusOK, payloadCity{
			payloadUser: createPayloadUser(c.Get("userName").(string)),
			City:        city,
		})
	}
}
