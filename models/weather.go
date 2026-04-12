// Package models structures the location and weather models
package models

type WeatherResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Days      []Day   `json:"days"`
	Current   Current `json:"currentConditions"`
}

type Day struct {
	DateTime   string  `json:"datetime"`
	TempMax    float64 `json:"tempmax"`
	TempMin    float64 `json:"tempmin"`
	Conditions string  `json:"conditions"`
	Humidity   float64 `json:"humidity"`
}

type Current struct {
	Temp       float64 `json:"temp"`
	FeelsLike  float64 `json:"feelslike"`
	Humidity   float64 `json:"humidity"`
	Conditions string  `json:"conditions"`
}

type Location struct {
	Location string `json:"location"`
}
