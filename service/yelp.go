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
	SearchByLocation(term string, location string) (model.SearchResponse, error)
	SearchByCoordinates(term string, latitude float64, longitude float64) (model.SearchResponse, error)
}

type yelpService struct {
	APIKey string
}

func NewYelpService(apiKey string) YelpService {
	return yelpService{
		APIKey: apiKey,
	}
}

func (svc yelpService) SearchByLocation(term string, location string) (model.SearchResponse, error) {
	searchURL := fmt.Sprintf(
		"https://api.yelp.com/v3/businesses/search?term=%s&location=%s",
		url.QueryEscape(term),
		url.QueryEscape(location),
	)

	return svc.search(searchURL)
}

func (svc yelpService) SearchByCoordinates(term string, latitude float64, longitude float64) (model.SearchResponse, error) {
	searchURL := fmt.Sprintf(
		"https://api.yelp.com/v3/businesses/search?term=%s&latitude=%f&longitude=%f",
		url.QueryEscape(term),
		latitude,
		longitude,
	)

	return svc.search(searchURL)
}

func (svc yelpService) search(url string) (model.SearchResponse, error) {
	res, err := doSearchRequest(url, svc.APIKey)
	if err != nil {
		return model.SearchResponse{}, err
	}

	defer res.Body.Close()

	log.Printf("[SearchByCoordinates] Response status: %s", res.Status)

	searchResponse, err := handleSearchResponse(res)

	return searchResponse, nil
}

func doSearchRequest(url string, key string) (*http.Response, error) {
	log.Printf("[doSearchRequest] Search request to Yelp: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("[doSearchRequest] failed to create GET request for /businesses/search bot, %s", err.Error())
	}
	addBearerToken(req, key)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do GET request to %s: %s", req.URL.RawPath, err.Error())
	}

	return res, nil
}

func addBearerToken(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
}

func handleSearchResponse(res *http.Response) (model.SearchResponse, error) {
	var searchResponse model.SearchResponse

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
