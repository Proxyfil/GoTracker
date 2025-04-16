package fdcnal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var apiEndpoint = "https://api.nal.usda.gov/fdc/v1/foods/"
var apiToken = "fEU6wcmaqSdJrMXeSPXnn19Lph66mtLsUCSEP67j"

// Struct to match the expected structure of the API response
type FoodSearchResponse struct {
	Foods []struct {
		Description string `json:"description"`
		FdcId       int    `json:"fdcId"`
	} `json:"foods"`
}

func GetFoodByName(foodName string) ([]string, error) {
	// Build the URL
	url := fmt.Sprintf("%ssearch?query=%s&api_key=%s", apiEndpoint, foodName, apiToken)

	// Make the request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Parse JSON
	var result FoodSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// Extract and return food descriptions
	var foods []string
	for _, food := range result.Foods {
		foods = append(foods, food.Description + " | ID : " + fmt.Sprint(food.FdcId))
	}

	return foods, nil
}