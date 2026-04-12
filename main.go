package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/KernelH132/weather-api/handler"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	r := http.NewServeMux()

	r.HandleFunc("/getweather", handler.GetWeather)

	fmt.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}

}
