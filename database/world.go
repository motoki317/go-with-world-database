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
	Code           string   `json:"code,omitempty"  db:"Code"`
	Name           string   `json:"name,omitempty"  db:"Name"`
	Continent      string   `json:"continent,omitempty"  db:"Continent"`
	Region         string   `json:"region,omitempty"  db:"Region"`
	SurfaceArea    float64  `json:"surface_area,omitempty"  db:"SurfaceArea"`
	IndepYear      *int     `json:"indep_year,omitempty"  db:"IndepYear"`
	Population     int      `json:"population,omitempty"  db:"Population"`
	LifeExpectancy *float64 `json:"life_expectancy,omitempty"  db:"LifeExpectancy"`
	GNP            *float64 `json:"GNP,omitempty"  db:"GNP"`
	GNPOld         *float64 `json:"GNP_old,omitempty"  db:"GNPOld"`
	LocalName      string   `json:"local_name,omitempty"  db:"LocalName"`
	GovernmentForm string   `json:"government_form,omitempty"  db:"GovernmentForm"`
	HeadOfState    *string  `json:"head_of_state,omitempty"  db:"HeadOfState"`
	Capital        *int     `json:"capital,omitempty"  db:"Capital"`
	Code2          string   `json:"code2,omitempty"  db:"Code2"`
}

type payloadUser struct {
	Username string `json:"username,omitempty"`
}

type payloadCountry struct {
	payloadUser
	Country Country `json:"country,omitempty"`
}

type payloadCountries struct {
	payloadUser
	Countries []Country `json:"countries,omitempty"`
}

type payloadCity struct {
	payloadUser
	City City `json:"city,omitempty"`
}

func createPayloadUser(c echo.Context) payloadUser {
	return payloadUser{
		Username: c.Get("userName").(string),
	}
}

// MakeRetrieveCountriesHandler Worldデータベースから全ての国データを取得します
func MakeRetrieveCountriesHandler(db *sqlx.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		countries := []Country{}
		db.Select(&countries, "SELECT * FROM country")

		if len(countries) == 0 {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, payloadCountries{
			payloadUser: createPayloadUser(c),
			Countries: countries,
		})
	}
}

// MakeRetrieveCountryHandler Worldデータベースから国情報を取得します
func MakeRetrieveCountryHandler(db *sqlx.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		countryName := c.Param("countryName")
		
		country := Country{}
		db.Get(&country, "SELECT * FROM country WHERE Name=?", countryName)
		if (country.Name == "") {
			return c.NoContent(http.StatusNotFound)
		}
		return c.JSON(http.StatusOK, payloadCountry{
			payloadUser: createPayloadUser(c),
			Country: country,
		})
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
			payloadUser: createPayloadUser(c),
			City:        city,
		})
	}
}
