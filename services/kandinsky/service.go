package kandinsky

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"time"

	"github.com/paavill/awesome-tagger-bot/domain/services"
)

const (
	headerNameKey    = "X-Key"
	headerNameSecret = "X-Secret"
)

func NewService(host, key, secret string) services.Kandinsky {

	headers := map[string]string{
		headerNameKey:    "Key " + key,
		headerNameSecret: "Secret " + secret,
	}

	return &knd{
		host:    host,
		key:     key,
		secret:  secret,
		headers: headers,
	}

}

type knd struct {
	host    string
	key     string
	secret  string
	headers map[string]string
}

func (k *knd) GenerateImage(query string) (*image.Image, error) {
	rawImages, err := k.runGeneration(query)
	if err != nil {
		return nil, fmt.Errorf("error while running generation due: " + err.Error())
	}

	if len(rawImages) != 1 {
		return nil, fmt.Errorf("error while getting images due: images count is not 1")
	}

	rawImage := rawImages[0]

	img, err := k.decodeImage(rawImage)
	if err != nil {
		return nil, fmt.Errorf("error while decoding image due: " + err.Error())
	}

	return img, nil
}

func (k *knd) decodeImage(rawImage string) (*image.Image, error) {
	decodedImage := make([]byte, len(rawImage))
	_, err := base64.RawStdEncoding.Decode(decodedImage, []byte(rawImage))
	if err != nil {
		return nil, fmt.Errorf("error while decoding base64 image due: " + err.Error())
	}

	parsedImage, _, err := image.Decode(bytes.NewReader(decodedImage))
	if err != nil {
		return nil, fmt.Errorf("error while decoding image due: " + err.Error())
	}

	return &parsedImage, nil
}

func (k *knd) runGeneration(query string) ([]string, error) {
	path := "/key/api/v1/text2image/run"

	req := k.prepareRequest(path, http.MethodPost)

	modelId, err := k.getModelId()
	if err != nil {
		return nil, fmt.Errorf("error while getting model id due: " + err.Error())
	}

	reqBody := generateRequest{
		ModelId: modelId,
		Params: generateParams{
			Type:      generate,
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
	if err != nil {
		return nil, fmt.Errorf("error while sending request due: " + err.Error())
	}
	defer response.Body.Close()

	rawResponse, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response body due: " + err.Error())
	}

	generateResponse := &generateResponse{}
	err = json.Unmarshal(rawResponse, generateResponse)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling response body due: " + err.Error())
	}

	log.Printf("generateResponse: %+v", generateResponse)

	if generateResponse.ModelStatus == "DISABLED_BY_QUEUE" {
		return nil, fmt.Errorf("model is disabled by queue")
	}

	uuid := generateResponse.Uuid
	images, err := k.checkGeneration(uuid)
	if err != nil {
		return nil, fmt.Errorf("error while checking generation due: " + err.Error())
	}
	return images, nil
}

func (k *knd) checkGeneration(uuid string) ([]string, error) {
	path := "/key/api/v1/text2image/status/" + uuid
	for retryCount := 0; retryCount < 120; retryCount++ {
		req := k.prepareRequest(k.host+path, http.MethodGet)
		resp, err := http.DefaultClient.Do(&req)
		if err != nil {
			return nil, fmt.Errorf("error while sending request due: " + err.Error())
		}
		defer resp.Body.Close()
		rawResponse, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error while reading response body due: " + err.Error())
		}
		checkGenerateStatusResponse := &checkGenerateStatusResponse{}
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

func (k *knd) getModelId() (int, error) {
	path := "/key/api/v1/models"
	req := k.prepareRequest(path, http.MethodGet)

	response, err := http.DefaultClient.Do(&req)
	if err != nil {
		return 0, fmt.Errorf("error while getting model id due: " + err.Error())
	}
	defer response.Body.Close()

	rawResponse, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, fmt.Errorf("error while reading response body due: " + err.Error())
	}

	models := &[]getModelRequest{}
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

func (k *knd) prepareRequest(path, method string) http.Request {
	url, err := url.Parse(k.host + path)
	if err != nil {
		log.Println("error while parsing url due: " + err.Error())
	}

	headers := http.Header{}
	for k, v := range k.headers {
		headers.Add(k, v)
	}

	req := http.Request{
		Method: method,
		URL:    url,
		Header: headers,
	}

	return req
}
