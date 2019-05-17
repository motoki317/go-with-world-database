package database

import (
	"fmt"
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

type payloadCities struct {
	payloadUser
	Cities []City `json:"cities,omitempty"`
}

func createPayloadUser(c echo.Context) payloadUser {
	return payloadUser{
		Username: c.Get("userName").(string),
	}
}

// MakeGetCountriesHandler Worldデータベースから全ての国データを取得します
func MakeGetCountriesHandler(db *sqlx.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		countries := []Country{}
		err := db.Select(&countries, "SELECT * FROM country")

		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		if len(countries) == 0 {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, payloadCountries{
			payloadUser: createPayloadUser(c),
			Countries: countries,
		})
	}
}

// MakeGetCountryHandler Worldデータベースから国情報を取得します
func MakeGetCountryHandler(db *sqlx.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		countryName := c.Param("countryName")
		
		country := Country{}
		err := db.Get(&country, "SELECT * FROM country WHERE Name=?", countryName)

		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		if (country.Name == "") {
			return c.NoContent(http.StatusNotFound)
		}
		return c.JSON(http.StatusOK, payloadCountry{
			payloadUser: createPayloadUser(c),
			Country: country,
		})
	}
}

// MakeGetCityHandler Worldデータベースから都市情報を取得します
func MakeGetCityHandler(db *sqlx.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		cityName := c.Param("cityName")

		city := City{}
		err := db.Get(&city, "SELECT * FROM city WHERE Name=?", cityName)

		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		if city.Name == "" {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusOK, payloadCity{
			payloadUser: createPayloadUser(c),
			City:        city,
		})
	}
}

// MakeGetCitiesInCountryHandler CountryNameからその国のCity一覧を取得します
func MakeGetCitiesInCountryHandler(db *sqlx.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		countryName := c.Param("countryName")

		cities := []City{}
		err := db.Select(&cities, "SELECT city.ID, city.Name, city.District, city.Population FROM city JOIN country ON country.Code = city.CountryCode WHERE country.Name=?", countryName)
		
		if err != nil {
			fmt.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
		if len(cities) == 0 {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusOK, payloadCities{
			payloadUser: createPayloadUser(c),
			Cities: cities,
		})
	}
}
