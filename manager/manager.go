package manager

import (
	"YandexEdaParser/db"
	"YandexEdaParser/structs"
)

type (
	YandexManager struct {
		Repository db.Repository
		Latitude   float64
		Longitude  float64
		Rating     float64
	}
)

func (y *YandexManager) Run(db structs.DataBase, loc structs.Location, rating float64) {
	y.Repository.Init(db)
	y.Longitude = loc.Longitude
	y.Latitude = loc.Latitude
	y.Rating = rating

}
