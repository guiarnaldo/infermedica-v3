package infermedica

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"strconv"
)

type SearchRes struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}
type LabTestsSearchRes struct {
	ID      string      `json:"id"`
	Label   string      `json:"label"`
	Results []LabResult `json:"results"`
}

type SearchType string

const (
	SearchTypeSymptom    SearchType = "symptom"
	SearchTypeRiskFactor SearchType = "risk_factor"
	SearchTypeLabTest    SearchType = "lab_test"
)

func (s SearchType) Ptr() *SearchType { return &s }
func (s SearchType) String() string   { return string(s) }

func (s *SearchType) IsValid() error {
	_, err := SearchTypeFromString(s.String())
	if err != nil {
		return err
	}
	return nil
}

func SearchTypeFromString(x string) (SearchType, error) {
	switch strings.ToLower(x) {
	case "symptom":
		return SearchTypeSymptom, nil
	case "risk_factor":
		return SearchTypeRiskFactor, nil
	case "lab_test":
		return SearchTypeLabTest, nil
	default:
		return "", fmt.Errorf("infermedica: unexpected value for search type: %q", x)
	}
}

func (a *App) Search(phrase string, sex Sex, maxResults int, st SearchType) (*[]SearchRes, error) {
	if sex.IsValid() != nil {
		return nil, fmt.Errorf("infermedica: Unexpected value for Sex")
	}
	if st.IsValid() != nil {
		return nil, fmt.Errorf("infermedica: Unexpected value for search type")
	}
	url := "search?phrase=" + url.QueryEscape(phrase) + "&sex=" + sex.String() + "&max_results=" + strconv.Itoa(maxResults) + "&type=" + st.String()
	req, err := a.prepareRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = checkResponse(res)

	if err != nil {
		return nil, err
	}
	var r []SearchRes
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
