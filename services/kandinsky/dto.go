package kandinsky

type generateType string

const (
	generate generateType = "GENERATE"
)

type generateParams struct {
	Type           generateType `json:"type"`
	NumImages      uint         `json:"numImages"`
	Width          uint         `json:"width"`
	Height         uint         `json:"height"`
	GenerateParams struct {
		Query string `json:"query"`
	} `json:"generateParams"`
}

type generateRequest struct {
	ModelId int            `json:"model_id"`
	Params  generateParams `json:"params"`
}

type generateResponseStatus string

const (
	GenerateResponseStatusInitial    generateResponseStatus = "INITIAL"
	GenerateResponseStatusProcessing generateResponseStatus = "PROCESSING"
	GenerateResponseStatusDone       generateResponseStatus = "DONE"
	GenerateResponseStatusFail       generateResponseStatus = "FAIL"
)

type generateResponse struct {
	Uuid        string                 `json:"uuid"`
	Status      generateResponseStatus `json:"status"`
	ModelStatus string                 `json:"model_status"`
}

type checkGenerateStatusResponse struct {
	Uuid             string                 `json:"uuid"`
	Status           generateResponseStatus `json:"status"`
	Images           []string               `json:"images"`
	ErrorDescription string                 `json:"errorDescription"`
	Censored         bool                   `json:"censored"`
}

type getModelRequest struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Version float64 `json:"version"`
	Type    string  `json:"type"`
}
