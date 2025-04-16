package fdcnal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var apiEndpoint = "https://api.nal.usda.gov/fdc/v1/"
var apiToken = "fEU6wcmaqSdJrMXeSPXnn19Lph66mtLsUCSEP67j"

// Struct to match the expected structure of the API response
type FoodSearchResponse struct {
	Foods []struct {
		Description string `json:"description"`
		FdcId       int    `json:"fdcId"`
	} `json:"foods"`
}

// Struct to match the expected structure of the API response for food details
type FoodDetailsResponse struct {
	Description string `json:"description"`
	FdcId       int    `json:"fdcId"`
	Nutrients   []struct {
		Id 	 int     `json:"id"`
		Amount float64 `json:"amount"`
	} `json:"foodNutrients"`
}

func GetFoodByName(foodName string) ([]string, error) {
	// Build the URL
	url := fmt.Sprintf("%sfoods/search?query=%s&api_key=%s", apiEndpoint, foodName, apiToken)

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

func GetFoodDetails(fdcId string) ([]string, error) {
	// Build the URL
	url := fmt.Sprintf("%sfood/%s?api_key=%s", apiEndpoint, fdcId, apiToken)

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

	fmt.Println(string(body))

	var foodDetails FoodDetailsResponse
	if err := json.Unmarshal(body, &foodDetails); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// Extract and return food details
	var details []string
	details = append(details, fmt.Sprintf("Description: %s", foodDetails.Description))
	details = append(details, fmt.Sprintf("FDC ID: %d", foodDetails.FdcId))
	details = append(details, "Nutrients:")

	for _, nutrient := range foodDetails.Nutrients {
		details = append(details, fmt.Sprintf("%d: %.2f", nutrient.Id, nutrient.Amount))
	}

	return details, nil
}