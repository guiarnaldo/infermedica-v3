package infermedica

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

// Check response status, returns an error with the status code if it is not 200
func checkResponse(res *http.Response) error {
	if res.StatusCode != http.StatusOK {
		var response Response
		json.NewDecoder(res.Body).Decode(&response)
		return fmt.Errorf("infermedica: %s: %s", res.Status, response.Message)
	}
	return nil
}

func (a *App) prepareRequest(method, url string, body interface{}) (*http.Request, error) {
	switch method {
	case "GET":
		return a.prepareGETRequest(url)
	case "POST":
		return a.preparePOSTRequest(url, body)
	}
	return nil, fmt.Errorf("infermedica: method not allowed")
}

func (a *App) addHeaders(req *http.Request) {
	req.Header.Add("App-Id", a.appID)
	req.Header.Add("App-Key", a.appKey)
	req.Header.Add("Content-Type", "application/json")
	if a.devMode {
		req.Header.Add("Dev-Mode", "true")
	}
	if a.model != "" {
		req.Header.Add("Model", a.model)
	}
	if a.interviewID != "" {
		req.Header.Add("Interview-Id", a.interviewID)
	}
}

func (a *App) prepareGETRequest(url string) (*http.Request, error) {
	baseURL := a.baseURL
	req, err := http.NewRequest("GET", baseURL+url, nil)
	if err != nil {
		return nil, err
	}
	a.addHeaders(req)
	return req, nil
}

func (a *App) preparePOSTRequest(url string, body interface{}) (*http.Request, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		return nil, err
	}
	baseURL := a.baseURL
	req, err := http.NewRequest("POST", baseURL+url, b)
	if err != nil {
		return nil, err
	}
	a.addHeaders(req)
	return req, nil
}
