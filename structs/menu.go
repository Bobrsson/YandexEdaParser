package structs

import (
	"time"
)

type ResponseMenu struct {
	Payload PayloadMenu `json:"payload"`
}

type PayloadMenu struct {
	Categories []Categories `json:"categories"`
}

type Categories struct {
	Restaurant Restaurant `json:"place"`
	Items      []Items    `json:"items"`
}

type Items struct {
	ID             int
	YandexId       int `json:"id"`
	RestaurantId   int
	CreateAt       time.Time
	UpdateAt       time.Time
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Price          float64 `yaml:"price"`
	Weight         string  `yaml:"weight"`
	Measure        Measure `yaml:"measure"`
	InternalRating float64
}

type Measure struct {
	Value    string `yaml:"value"`
	ValueInt int
	Measure  string `yaml:"measure_unit"`
}
