package get_image

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

var (
	kandinskyApiHost   = "https://api-key.fusionbrain.ai"
	kandinskyTypesHost = "https://cdn.fusionbrain.ai/static/styles/key"
	apiKey             = "Key 0FE09D3D0FD367E48EBD7A6A6A66F541"
	secretKey          = "Secret 5B9B8424D6046478E11815B029EDD7EC"
	keyHeader          = "X-Key"
	secretHeader       = "X-Secret"
	headers            = map[string]string{
		keyHeader:    apiKey,
		secretHeader: secretKey,
	}
)

type GenerateType string

const (
	Generate GenerateType = "GENERATE"
)

type GenerateParams struct {
	Type           GenerateType `json:"type"`
	NumImages      uint         `json:"numImages"`
	Width          uint         `json:"width"`
	Height         uint         `json:"height"`
	GenerateParams struct {
		Query string `json:"query"`
	} `json:"generateParams"`
}

type GenerateRequest struct {
	ModelId string         `json:"model_id"`
	Params  GenerateParams `json:"params"`
}

type GenerateResponse struct {
	Uuid        string `json:"uuid"`
	Status      string `json:"status"`
	ModelStatus string `json:"model_status"`
}

type CheckGenerateStatusResponse struct {
	Uuid             string   `json:"uuid"`
	Status           string   `json:"status"`
	Images           []string `json:"images"`
	ErrorDescription string   `json:"errorDescription"`
	Censored         bool     `json:"censored"`
}

type GetModel struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Type    string `json:"type"`
}

func Run(query string) ([]string, error) {
	err := generate()
	if err != nil {
		return nil, err
	}
}

func generate(query string) error {
	path := "/key/api/v1/text2image/run"

	req := prepareRequest(kandinskyApiHost+path, http.MethodPost, headers)

	modelId, err := getModelId()
	if err != nil {
		return fmt.Errorf("error while getting model id due: " + err.Error())
	}

	reqBody := GenerateRequest{
		ModelId: modelId,
		Params: GenerateParams{
			Type:      Generate,
			NumImages: 1,
			Width:     1024,
			Height:    1024,
			GenerateParams: struct {
				Query string `json:"query"`
			}{
				Query: query,
			},
		},
	}

	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("error while marshaling request body due: " + err.Error())
	}

	req.Body = io.NopCloser(bytes.NewBuffer(reqBodyJson))

	//response, err := http.DefaultClient.Do()
}

func checkGeneration(uuid string) {
	path := "/key/api/v1/text2image/status/" + uuid

}

func getModelId() (string, error) {
	path := "key/api/v1/models"
	req := prepareRequest(kandinskyApiHost+path, http.MethodGet, headers)

	response, err := http.DefaultClient.Do(&req)
	defer response.Body.Close()
	if err != nil {
		return "", fmt.Errorf("error while getting model id due: " + err.Error())
	}

	rawResponse, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading response body due: " + err.Error())
	}

	models := []GetModel{}
	err = json.Unmarshal(rawResponse, models)
	if err != nil {
		return "", fmt.Errorf("error while unmarshalling response body due: " + err.Error())
	}

	if len(models) == 0 {
		return "", fmt.Errorf("error while unmarshalling response body due: no models found")
	}

	return models[0].Id, nil
}

func prepareRequest(urlRaw string, method string, headersRaw map[string]string) http.Request {
	url, err := url.Parse(urlRaw)
	if err != nil {
		log.Println("error while parsing url due: " + err.Error())
	}

	headers := http.Header{}
	for k, v := range headersRaw {
		headers.Add(k, v)
	}

	req := http.Request{
		Method: method,
		URL:    url,
		Header: headers,
	}

	return req
}
