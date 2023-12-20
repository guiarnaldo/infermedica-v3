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

type SearchReq struct {
	Phrase     string     `json:"phrase"`
	Sex        Sex        `json:"sex"`
	Age        Age        `json:"age"`
	MaxResults int        `json:"max_results"`
	Types      SearchType `json:"types"`
}

type SearchRes struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type SearchType string

const (
	SearchTypeSymptom    SearchType = "symptom"
	SearchTypeRiskFactor SearchType = "risk_factor"
	SearchTypeLabTest    SearchType = "lab_test"
	SearchTypeCondition  SearchType = "condition"
)

func (s *SearchType) IsValid() error {
	_, err := SearchTypeFromString(string(*s))
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

// Search returns a list of observations matching the given phrase.
func (a *App) Search(sq SearchReq) (*[]SearchRes, error) {
	if sq.Sex.IsValid() != nil {
		return nil, fmt.Errorf("infermedica: Unexpected value for Sex")
	}
	if sq.Types.IsValid() != nil {
		return nil, fmt.Errorf("infermedica: Unexpected value for search type")
	}
	if sq.MaxResults <= 0 {
		return nil, fmt.Errorf("infermedica: MaxResult can not be zero or less")
	}
	url := "search?phrase=" + url.QueryEscape(sq.Phrase) + "&sex=" + string(sq.Sex) +
		"&max_results=" + strconv.Itoa(sq.MaxResults) + "&types=" + string(sq.Types) + "&age.value=" + strconv.Itoa(sq.Age.Value) + "&age.unit=" + string(sq.Age.Unit)
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
