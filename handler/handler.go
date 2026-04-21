// Package handler sends the weather data to the client from the location provided
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/KernelH132/weather-api/models"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func GetWeather(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var location models.Location

	// Method Check
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/getweather" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	// Decode Input
	err := json.NewDecoder(r.Body).Decode(&location)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("weather:%s", location.Location)

	cachedData, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "HIT")
		fmt.Fprint(w, cachedData)
		return
	}

	// Get API Key
	apiKey := os.Getenv("WEATHER_KEY")
	if apiKey == "" {
		fmt.Println("Error: WEATHER_KEY environment variable is not set.")
		return
	}

	// Get Weather
	dateStr := time.Now().Format("2006-01-02")
	url := fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s/%s?key=%s&unitGroup=metric", location.Location, dateStr, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Weather service unreachable", http.StatusServiceUnavailable)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API error: status %d\n", resp.StatusCode)
		return
	}

	defer resp.Body.Close()

	var weather models.WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		http.Error(w, "Failed to parse weather data", http.StatusInternalServerError)
		return
	}

	// Combine and Send
	finalResult := struct {
		UserLocation models.Location        `json:"user_request"`
		Data         models.WeatherResponse `json:"weather_data"`
	}{
		UserLocation: location,
		Data:         weather,
	}

	jsonData, err := json.Marshal(finalResult)
	if err != nil {
		http.Error(w, "Something went wrong", 500)
		return
	}

	rdb.Set(ctx, cacheKey, jsonData, 20*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Cache", "MISS")
	w.Write(jsonData)
}
