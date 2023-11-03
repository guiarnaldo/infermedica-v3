package infermedica

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LabTestsRes struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	CommonName string      `json:"common_name"`
	Category   string      `json:"category"`
	Results    []LabResult `json:"results"`
}

type LabResult struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func (a *App) LabTests() (*[]LabTestsRes, error) {
	req, err := a.prepareRequest("GET", "lab_tests", nil)
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
	var r []LabTestsRes
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (a *App) LabTestsIDMap() (*map[string]LabTestsRes, error) {
	r, err := a.LabTests()
	if err != nil {
		return nil, err
	}
	rmap := make(map[string]LabTestsRes)
	for _, sr := range *r {
		rmap[sr.ID] = sr
	}
	return &rmap, nil
}

func (a *App) LabTestByID(id string) (*LabTestsRes, error) {
	req, err := a.prepareRequest("GET", "lab_tests/"+id, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var r LabTestsRes
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

type LabTestsRecommendRes struct {
	Recommended []LabTestsRecommendation `json:"recommended"`
	Obligatory  []LabTestsRecommendation `json:"obligatory"`
}
type LabTestsRecommendation struct {
	PanelID  string       `json:"panel_id"`
	Name     string       `json:"name"`
	Position int          `json:"position"`
	LabTests []LabTestsID `json:"lab_tests"`
}
type LabTestsID struct {
	ID string `json:"id"`
}

// Recommend is a func to request lab test recommendations for given data
func (a *App) LabTestsRecommend(dr ObservationReq) (*LabTestsRecommendRes, error) {
	if dr.Sex.IsValid() != nil {
		return nil, fmt.Errorf("infermedica: Unexpected value for Sex")
	}
	req, err := a.prepareRequest("POST", "lab_tests/recommend", dr)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var r LabTestsRecommendRes
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
