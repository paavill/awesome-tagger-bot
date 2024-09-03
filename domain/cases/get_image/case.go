package get_image

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"time"
)

var (
	kandinskyApiHost   = "https://api-key.fusionbrain.ai"
	kandinskyTypesHost = "https://cdn.fusionbrain.ai/static/styles/key"
	apiKey             = "Key"
	secretKey          = "Secret"
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
	ModelId int            `json:"model_id"`
	Params  GenerateParams `json:"params"`
}

type GenerateResponseStatus string

const (
	GenerateResponseStatusInitial    GenerateResponseStatus = "INITIAL"
	GenerateResponseStatusProcessing GenerateResponseStatus = "PROCESSING"
	GenerateResponseStatusDone       GenerateResponseStatus = "DONE"
	GenerateResponseStatusFail       GenerateResponseStatus = "FAIL"
)

type GenerateResponse struct {
	Uuid        string                 `json:"uuid"`
	Status      GenerateResponseStatus `json:"status"`
	ModelStatus string                 `json:"model_status"`
}

type CheckGenerateStatusResponse struct {
	Uuid             string                 `json:"uuid"`
	Status           GenerateResponseStatus `json:"status"`
	Images           []string               `json:"images"`
	ErrorDescription string                 `json:"errorDescription"`
	Censored         bool                   `json:"censored"`
}

type GetModelRequest struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Version float64 `json:"version"`
	Type    string  `json:"type"`
}

func Run(query string) ([]string, error) {
	images, err := generate(query)
	if err != nil {
		return nil, err
	}

	if len(images) != 1 {
		return nil, fmt.Errorf("expected 1 image, got %d", len(images))
	}

	base64Image := []byte(images[0])

	rawImage := make([]byte, len(base64Image))
	n, err := base64.RawStdEncoding.Decode(rawImage, base64Image)
	log.Println("decoded bytes:", n)
	if err != nil {
		return nil, fmt.Errorf("error while decoding base64 image due: " + err.Error())
	}

	img, _, err := image.Decode(bytes.NewReader(rawImage))
	if err != nil {
		return nil, fmt.Errorf("error while decoding image due: " + err.Error())
	}

	f, err := os.Create("image.jpg")
	if err != nil {
		return nil, fmt.Errorf("error while creating file due: " + err.Error())
	}
	defer f.Close()

	err = jpeg.Encode(f, img, nil)
	if err != nil {
		return nil, fmt.Errorf("error while encoding image due: " + err.Error())
	}

	return nil, nil
}

func generate(query string) ([]string, error) {
	path := "/key/api/v1/text2image/run"

	req := prepareRequest(kandinskyApiHost+path, http.MethodPost, headers)

	modelId, err := getModelId()
	if err != nil {
		return nil, fmt.Errorf("error while getting model id due: " + err.Error())
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

	reqBodyJsonParams, err := json.Marshal(reqBody.Params)
	if err != nil {
		return nil, fmt.Errorf("error while marshaling request body due: " + err.Error())
	}

	buffer := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(buffer)

	err = multipartWriter.WriteField("model_id", fmt.Sprint(reqBody.ModelId))
	if err != nil {
		return nil, fmt.Errorf("error while creating form field due: " + err.Error())
	}

	mimeHeader := textproto.MIMEHeader{}
	mimeHeader.Add("Content-Type", "application/json")
	mimeHeader.Add("Content-Disposition", `form-data; name="params"`)
	paramsWriter, err := multipartWriter.CreatePart(mimeHeader)
	if err != nil {
		return nil, fmt.Errorf("error while creating form field due: " + err.Error())
	}
	_, err = paramsWriter.Write(reqBodyJsonParams)
	if err != nil {
		return nil, fmt.Errorf("error while writing to form field due: " + err.Error())
	}

	err = multipartWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("error while closing multipart writer due: " + err.Error())
	}

	req.Body = io.NopCloser(buffer)

	contentType := multipartWriter.FormDataContentType()
	req.Header.Set("Content-Type", contentType)

	response, err := http.DefaultClient.Do(&req)
	defer response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error while sending request due: " + err.Error())
	}

	rawResponse, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response body due: " + err.Error())
	}

	generateResponse := &GenerateResponse{}
	err = json.Unmarshal(rawResponse, generateResponse)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling response body due: " + err.Error())
	}

	log.Printf("generateResponse: %+v", generateResponse)

	if generateResponse.ModelStatus == "DISABLED_BY_QUEUE" {
		return nil, fmt.Errorf("model is disabled by queue")
	}

	uuid := generateResponse.Uuid
	images, err := checkGeneration(uuid)
	if err != nil {
		return nil, fmt.Errorf("error while checking generation due: " + err.Error())
	}
	return images, nil
}

func checkGeneration(uuid string) ([]string, error) {
	path := "/key/api/v1/text2image/status/" + uuid
	for retryCount := 0; retryCount < 120; retryCount++ {
		req := prepareRequest(kandinskyApiHost+path, http.MethodGet, headers)
		resp, err := http.DefaultClient.Do(&req)
		defer resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("error while sending request due: " + err.Error())
		}
		rawResponse, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error while reading response body due: " + err.Error())
		}
		checkGenerateStatusResponse := &CheckGenerateStatusResponse{}
		err = json.Unmarshal(rawResponse, checkGenerateStatusResponse)
		if err != nil {
			return nil, fmt.Errorf("error while unmarshalling response body due: " + err.Error())
		}
		if checkGenerateStatusResponse.Status == GenerateResponseStatusDone {
			return checkGenerateStatusResponse.Images, nil
		}
		if checkGenerateStatusResponse.Status == GenerateResponseStatusFail {
			return nil, fmt.Errorf("error while generating image due: " + checkGenerateStatusResponse.ErrorDescription)
		}
		time.Sleep(1 * time.Second)
	}
	return nil, fmt.Errorf("error while checking generation status")
}

func getModelId() (int, error) {
	path := "/key/api/v1/models"
	req := prepareRequest(kandinskyApiHost+path, http.MethodGet, headers)

	response, err := http.DefaultClient.Do(&req)
	defer response.Body.Close()
	if err != nil {
		return 0, fmt.Errorf("error while getting model id due: " + err.Error())
	}

	rawResponse, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, fmt.Errorf("error while reading response body due: " + err.Error())
	}

	models := &[]GetModelRequest{}
	err = json.Unmarshal(rawResponse, models)
	if err != nil {
		return 0, fmt.Errorf("error while unmarshalling response body due: " + err.Error())
	}

	modelValue := *models

	if len(modelValue) == 0 {
		return 0, fmt.Errorf("error while unmarshalling response body due: no models found")
	}

	return modelValue[0].Id, nil
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
