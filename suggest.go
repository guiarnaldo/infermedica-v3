package infermedica

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SuggestReq is a struct to request suggestions
type SuggestReq struct {
	Sex           Sex           `json:"sex"`
	Age           Age           `json:"age"`
	Evidences     []Evidence    `json:"evidence"`
	SuggestMethod SuggestMethod `json:"suggest_method"`
	Extras        SuggestExtras `json:"extras"`
}

type SuggestExtras []struct {
	EnableExplanations bool `json:"enable_explanations"` // This functionality helps users to better understand the purpose of a question. It expands the question with two additional fields: explication and instruction
}

// SuggestRes is a response struct for suggest
type SuggestRes struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	CommonName  string   `json:"common_name"`
	Explication string   `json:"explication"` // Enabled only when EnableExplanations is true
	Instruction []string `json:"instruction"` // Enabled only when EnableExplanations is true
}

type SuggestMethod string

const (
	SuggestMethodSymptoms                 SuggestMethod = "symptoms"                    //Similar symptoms (default)
	SuggestMethodRiskFactors              SuggestMethod = "risk_factors"                //Relevant risk factors. This method was deprecated with the release of API 3.5
	SuggestMethoddemoGraphicRiskFactors   SuggestMethod = "demographic_risk_factors"    //Demographic risk factors
	SuggestMethodEvidenceBasedRiskFactors SuggestMethod = "evidence_based_risk_factors" // Evidence-based risk factors
	SuggestMethodRedFlags                 SuggestMethod = "red_flags"                   // Red flags
)

// Suggest is a func to request suggestions
func (a *App) Suggest(sr SuggestReq) (*[]SuggestRes, error) {
	if sr.Sex.IsValid() != nil {
		return nil, fmt.Errorf("infermedica: Unexpected value for Sex")
	}
	req, err := a.prepareRequest("POST", "suggest", sr)
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

	// Check response
	err = checkResponse(res)
	if err != nil {
		return nil, err
	}

	var r []SuggestRes
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
