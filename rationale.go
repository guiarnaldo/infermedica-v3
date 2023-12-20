package infermedica

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RationaleReq struct {
	Sex         Sex                 `json:"sex"`
	Age         Age                 `json:"age"`
	EvaluatedAt string              `json:"evaluated_at,omitempty"`
	Evidences   []Evidence          `json:"evidence,omitempty"`
	Extras      *RationaleReqExtras `json:"extras,omitempty"`
}

type RationaleReqExtras struct {
	EnableSymptomDuration bool `json:"enable_symptom_duration,omitempty"` // This flag enables questions of the type duration which contain a new field EvidenceID
}

type RationaleRes struct {
	Type              RationaleType       `json:"type"`
	ObservationParams []ObservationParams `json:"observation_params"`
	ConditionParams   []ConditionParams   `json:"condition_params"`
}
type ObservationParams struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CommonName string `json:"common_name"`
}
type ConditionParams struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CommonName string `json:"common_name"`
}

type RationaleType string

const (
	RationaleTypeR0 RationaleType = "r0" // I'm asking this question because a negative response reduces the probability of condition_params and other conditions
	RationaleTypeR1 RationaleType = "r1" // I'm asking this question because observation_params might be related to one or more considered conditions
	RationaleTypeR2 RationaleType = "r2" // I'm asking this question because a negative response reduces the probability of condition_params and other conditions
	RationaleTypeR3 RationaleType = "r3" // I'm asking this question to either rule in or out conditions such as condition_params
	RationaleTypeR4 RationaleType = "r4" // I'm asking this question because observation_params might be one of the causes of your symptoms
	RationaleTypeR5 RationaleType = "r5" // I'm asking this question because I want to know if you suffered any recent injuries
	RationaleTypeR6 RationaleType = "r6" // I'm asking this question to learn more about your observation_params
)

// Rationale returns the rationale behind the questions that are asked by the system
func (a *App) Rationale(sr RationaleReq) (*[]RationaleRes, error) {
	if sr.Sex.IsValid() != nil {
		return nil, fmt.Errorf("infermedica: Unexpected value for Sex")
	}
	req, err := a.prepareRequest("POST", "Rationale", sr)
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

	var r []RationaleRes
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
