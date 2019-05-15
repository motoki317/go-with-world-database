package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type City struct {
	ID 			int	   `json:"id,omitempty" db:"ID"`
	Name        string `json:"name,omitempty"  db:"Name"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

type Country struct {
	Name string `json:"name,omitempty" db:"Name"`
	Population int `json:"population,omitempty" db:"Population"`
}

func main() {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}

	fmt.Println("Connected!")

	selected := "Tokyo"
	if len(os.Args) > 1 {
		selected = ""
		for i, arg := range os.Args {
			if i != 0 {
				selected += arg + " "
			}
		}
		runes := []rune(selected)
		selected = string(runes[:len(runes)-1])
	}

	city := City{}
	db.Get(&city, "SELECT * FROM city WHERE Name='" + selected + "'")

	fmt.Println(city)
	country := Country{}
	err = db.Get(&country, "SELECT Name, Population FROM country WHERE Code='" + city.CountryCode + "'")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(country)

	fmt.Println(country)
	fmt.Printf("%sの人口は%d人です\n", selected, city.Population)
	fmt.Printf("%sの人口のうち%v%%を占めています\n", country.Name, (float64(city.Population) / float64(country.Population) * 100))
}