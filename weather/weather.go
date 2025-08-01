package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

const baseURL = "http://api.openweathermap.org/data/2.5/weather"

var icons = map[string]string{
	"01d":     " ",
	"01n":     " ",
	"02d":     " ",
	"03d":     " ",
	"03n":     " ",
	"04d":     " ",
	"04n":     " ",
	"02n":     "  ",
	"09":      " ",
	"10d":     " ",
	"10n":     " ",
	"10n 11n": " ",
	"10d 11d": " ",
	"11":      "",
	"13d":     " ",
	"13n":     " ",
	"50d":     " ",
	"50n":     " ",
}

type WeatherData struct {
	Weather []struct {
		Icon string `json:"icon"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type OutputData struct {
	Text  string `json:"text,omitempty"`
	Error string `json:"error,omitempty"`
}

func checkInternetConnection() bool {
	resp, err := http.Get("https://google.com")
	if err != nil {
		data := ErrorResponse{Error: "No Internet!"}
		jsonData, _ := json.Marshal(data.Error)
		fmt.Println(string(jsonData))
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

func getTemperature(city string) (*WeatherData, error) {
	if !checkInternetConnection() {
		return nil, fmt.Errorf("No Internet!")
	}

	apiKey := os.Getenv("WEATHER_API")
	if apiKey == "" {
		return nil, fmt.Errorf("WEATHER_API environment variable not set")
	}

	url := fmt.Sprintf("%s?q=%s&units=metric&APPID=%s", baseURL, city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var weatherData WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return &weatherData, nil
}

func displayHelp() {
	fmt.Println("Usage: weatherInfo [Options]")
	fmt.Println("Options:")
	fmt.Println("  -city <NAME> \t Name of the city of which to get weather information (required).")
	fmt.Println("  -env  <PATH> \t Path to the .env file (optional, defaults to ~/.env).")
	fmt.Println("\t\t Make sure it contains WEATHER_API=<your_api_key>")
	fmt.Println("  -help \t Display this help message.")
	os.Exit(0)
}

func main() {
	var cityFlag = flag.String("city", "", "City name to get the weather for")
	var envFlag = flag.String("env", "", "Path to the .env file (optional)")
	var helpFlag = flag.Bool("help", false, "Display help message")
	flag.Parse()

	if len(os.Args) < 2 {
		displayHelp()
	}

	if envFlag != nil && *envFlag != "" {
		// Load custom .env file if specified
		err := godotenv.Load(*envFlag)
		if err != nil {
			log.Fatalf("failed to load %s: %v", *envFlag, err)
		}
	} else {
		// Load default .env file from home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("failed to get home directory: %v", err)
		}
		envPath := filepath.Join(homeDir, ".env")
		err = godotenv.Load(envPath)
		if err != nil {
			log.Printf("failed to load %s: %v", envPath, err)
		}
	}

	if *helpFlag {
		displayHelp()
	}

	weatherData, err := getTemperature(*cityFlag)
	if err != nil {
		log.Fatal(err)
	}

	if len(weatherData.Weather) == 0 {
		log.Fatal("No weather data received")
	}

	code := weatherData.Weather[0].Icon
	icon, exists := icons[code]
	if !exists {
		icon = " " // Default icon if not found
	}

	temp := int(math.Ceil(weatherData.Main.Temp))

	data := OutputData{
		Text: fmt.Sprintf("%s %d°C", icon, temp),
	}

	fmt.Println(data.Text)
}
