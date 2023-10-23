package infermedica

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type App struct {
	baseURL     string
	appID       string
	appKey      string
	model       string
	interviewID string
}

func NewApp(id, key, model, interviewID string) App {
	return App{
		baseURL:     "https://api.infermedica.com/v3/",
		appID:       id,
		appKey:      key,
		model:       model,
		interviewID: interviewID,
	}
}

func (a App) prepareRequest(method, url string, body interface{}) (*http.Request, error) {
	switch method {
	case "GET":
		return a.prepareGETRequest(url, body)
	case "POST":
		return a.preparePOSTRequest(url, body)
	}
	return nil, errors.New("Method not allowed")
}

func (a App) addHeaders(req *http.Request) {
	req.Header.Add("App-Id", a.appID)
	req.Header.Add("App-Key", a.appKey)
	req.Header.Add("Content-Type", "application/json")
	if a.model != "" {
		req.Header.Add("Model", a.model)
	}
	if a.interviewID != "" {
		req.Header.Add("Interview-Id", a.interviewID)
	}
}

func (a App) prepareGETRequest(url string, body interface{}) (*http.Request, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		return nil, err
	}
	baseURL := a.baseURL
	req, err := http.NewRequest("GET", baseURL+url, b)
	if err != nil {
		return nil, err
	}
	a.addHeaders(req)
	return req, nil
}

func (a App) preparePOSTRequest(url string, body interface{}) (*http.Request, error) {
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

type Sex string

const (
	SexMale   Sex = "male"
	SexFemale Sex = "female"
)

func (s Sex) Ptr() *Sex      { return &s }
func (s Sex) String() string { return string(s) }

func (s *Sex) IsValid() error {
	_, err := SexFromString(s.String())
	if err != nil {
		return err
	}
	return nil
}

func SexFromString(x string) (Sex, error) {
	switch strings.ToLower(x) {
	case "male":
		return SexMale, nil
	case "female":
		return SexFemale, nil
	default:
		return "", fmt.Errorf("unexpected value for Sex: %q", x)
	}
}

type SexFilter string

const (
	SexFilterBoth   SexFilter = "both"
	SexFilterMale   SexFilter = "male"
	SexFilterFemale SexFilter = "female"
)

func (s SexFilter) Ptr() *SexFilter { return &s }
func (s SexFilter) String() string  { return string(s) }

func (s *SexFilter) IsValid() error {
	_, err := SexFilterFromString(s.String())
	if err != nil {
		return err
	}
	return nil
}

func SexFilterFromString(x string) (SexFilter, error) {
	switch strings.ToLower(x) {
	case "both":
		return SexFilterBoth, nil
	case "male":
		return SexFilterMale, nil
	case "female":
		return SexFilterFemale, nil
	default:
		return "", fmt.Errorf("unexpected value for SexFilter: %q", x)
	}
}

type EvidenceChoiceID string

const (
	EvidenceChoiceIDPresent EvidenceChoiceID = "present"
	EvidenceChoiceIDAbsent  EvidenceChoiceID = "absent"
	EvidenceChoiceIDUnknown EvidenceChoiceID = "unknown"
)

func (ecID EvidenceChoiceID) Ptr() *EvidenceChoiceID { return &ecID }
func (ecID EvidenceChoiceID) String() string         { return string(ecID) }

func (ecID EvidenceChoiceID) IsValid() error {
	_, err := EvidenceChoiceIDFromString(ecID.String())
	if err != nil {
		return err
	}
	return nil
}

func EvidenceChoiceIDFromString(x string) (EvidenceChoiceID, error) {
	switch strings.ToLower(x) {
	case "present":
		return EvidenceChoiceIDPresent, nil
	case "absent":
		return EvidenceChoiceIDAbsent, nil
	case "unknown":
		return EvidenceChoiceIDUnknown, nil
	default:
		return "", fmt.Errorf("unexpected value for evidence choice id: %q", x)
	}
}

// Contains source valid types 
type EvidenceSource string

const (
	EvidenceSourceInitial EvidenceSource = "initial"
	EvidenceSourceSuggest  EvidenceSource = "suggest"
	EvidenceSourcePredefined EvidenceSource = "predefined"
	EvidenceSourceRedFlags EvidenceSource = "red_flags"
)

type Evidence struct {
	ID       string           `json:"id"`
	ChoiceID EvidenceChoiceID `json:"choice_id"`
	Source   EvidenceSource   `json:"source"`
}

type Age struct{
	Value int `json:"value"`
}