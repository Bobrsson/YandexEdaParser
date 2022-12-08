package structs

import (
	"time"
)

type ResponseRestaurant struct {
	Payload PayloadRestaurant `json:"payload"`
}

type PayloadRestaurant struct {
	Title       string   `json:"title"`
	FoundPlaces []Places `json:"foundPlaces"`
}

type Places struct {
	Restaurant Restaurant `json:"place"`
}

type Restaurant struct {
	ID                int
	YandexId          int `json:"id"`
	CreateAt          time.Time
	UpdateAt          time.Time
	Name              string  `json:"name"`
	Slug              string  `json:"slug"`
	Rating            float64 `yaml:"rating"`
	MinimalOrderPrice float64 `yaml:"MinimalOrderPrice"`
	InternalRating    float64 `yaml:"InternalRating"`
}
