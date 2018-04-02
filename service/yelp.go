package service

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/zachvanuum/FoodHelperBot/model"
	"github.com/zachvanuum/FoodHelperBot/util"
)

type YelpService interface {
	Search(term string) (model.SearchResponse, error)
}

type yelpService struct {
	APIKey string
}

func NewYelpService(apiKey string) YelpService {
	return yelpService{
		APIKey: apiKey,
	}
}

func (svc yelpService) Search(term string) (model.SearchResponse, error) {
	var searchResponse model.SearchResponse

	location := "Phoenix, AZ"
	searchURL := fmt.Sprintf(
		"https://api.yelp.com/v3/businesses/search?term=%s&location=%s",
		url.QueryEscape(term),
		url.QueryEscape(location),
	)

	log.Printf("[Search] Search request to Yelp: %s", searchURL)

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return searchResponse, fmt.Errorf("[Search] failed to create GET request for /businesses/search bot, %s", err.Error())
	}
	addBearerToken(req, svc.APIKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return searchResponse, fmt.Errorf("failed to do GET request to %s: %s", req.URL.RawPath, err.Error())
	}

	defer res.Body.Close()

	log.Printf("[Search] Response status: %s", res.Status)

	if err := util.UnmarshalBody(res.Body, &searchResponse); err != nil {
		return searchResponse, fmt.Errorf("failed to marshall search response to struct: %s", err.Error())
	}

	log.Printf("[Search] Search got %d results", searchResponse.Total)

	if searchResponse.Total == 0 {
		log.Printf("[Search] Checking for error response")

		var errorResponse model.ErrorResponseWrapper
		if err := util.UnmarshalBody(res.Body, &errorResponse); err != nil {
			return searchResponse, fmt.Errorf("failed to marshall error response to struct: %s", err.Error())
		}

		log.Printf("[Search] Error response - Code: %s, Description: %s", errorResponse.Error.Code, errorResponse.Error.Description)
	}

	return searchResponse, nil
}

func addBearerToken(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
}
