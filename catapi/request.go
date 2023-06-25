package catapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

var apiEnv = "API_KEY"

type Request struct {
	Breed string
}

func (r Request) Execute() ([]Response, error) {
	apiKey := os.Getenv(apiEnv)
	if len(apiKey) == 0 {
		return nil, errors.New(fmt.Sprintf("%s not set", apiEnv))
	}

	req, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/images/search", nil)
	if err != nil {
		return nil, err
	}
	req.Header = make(map[string][]string)
	req.Header.Set("x-api-key", apiKey)
	if len(r.Breed) != 0 {
		qry := req.URL.Query()
		qry.Add("breed_ids", r.Breed)
		req.URL.RawQuery = qry.Encode()
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var cats []Response
	err = json.NewDecoder(resp.Body).Decode(&cats)
	return cats, err
}

type BreedRequest struct {
}

func (r BreedRequest) Execute() ([]Breed, error) {
	resp, err := http.Get("https://api.thecatapi.com/v1/breeds")
	if err != nil {
		return nil, err
	}

	var breeds []Breed
	err = json.NewDecoder(resp.Body).Decode(&breeds)
	return breeds, err
}
