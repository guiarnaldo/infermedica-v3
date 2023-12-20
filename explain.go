package infermedica

import (
	"encoding/json"
	"net/http"
	"time"
)

type ExplainReq struct {
	Sex         Sex               `json:"sex"`
	Age         Age               `json:"age"`
	EvaluatedAt string            `json:"evaluated_at,omitempty"`
	Evidences   *[]Evidence       `json:"evidence,omitempty"`
	Extras      *ExplainReqExtras `json:"extras,omitempty"`
	Target      string            `json:"target"`
}

type ExplainReqExtras struct {
	EnableSymptomDuration bool `json:"enable_symptom_duration,omitempty"` // This flag enables questions of the type duration which contain a new field EvidenceID
}

type Observations struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CommonName string `json:"common_name"`
}

type ExplainRes struct {
	SupportingEvidence  []Observations `json:"supporting_evidence"`
	ConflictingEvidence []Observations `json:"conflicting_evidence"`
	UnconfirmedEvidence []Observations `json:"unconfirmed_evidence"`
}

// Explains which evidence impacts the probability of a selected condition appearing in the ranking
func (a *App) Explain(er ExplainReq) (*ExplainRes, error) {
	req, err := a.prepareRequest("POST", "explain", er)
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
	var r ExplainRes
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
