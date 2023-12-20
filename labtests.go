package infermedica

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type LabTestsReq struct {
	Sex         Sex             `json:"sex"`
	Age         Age             `json:"age"`
	EvaluatedAt string          `json:"evaluated_at,omitempty"`
	Evidences   []Evidence      `json:"evidence,omitempty"`
	Extras      *LabTestsExtras `json:"extras,omitempty"`
}

type LabTestsExtras struct {
	EnableSymptomDuration bool `json:"enable_symptom_duration,omitempty"` // This flag enables questions of the type duration which contain a new field EvidenceID
}

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

func (a *App) LabTests(age Age, enableTriage3 bool) (*[]LabTestsRes, error) {
	req, err := a.prepareRequest("GET", "lab_tests?age.value="+strconv.Itoa(age.Value)+"&age.unit"+string(age.Unit)+"&enableTriage3="+strconv.FormatBool(enableTriage3), nil)
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

func (a *App) LabTestByID(id string, age Age, enableTriage3 bool) (*LabTestsRes, error) {
	req, err := a.prepareRequest("GET", "lab_tests/"+id+"?age.value="+strconv.Itoa(age.Value)+"&age.unit"+string(age.Unit)+"&enableTriage3="+strconv.FormatBool(enableTriage3), nil)
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
func (a *App) LabTestsRecommend(dr LabTestsReq) (*LabTestsRecommendRes, error) {
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
