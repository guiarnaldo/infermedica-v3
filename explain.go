package infermedica

import (
	"encoding/json"
	"net/http"
	"time"
)

type ExplainReq struct {
	ObservationReq
	Target string `json:"target"` // ID of the condition that you want explained
}

type ExplainRes struct {
	SupportingEvidence  []Evidence `json:"supporting_evidence"`
	ConflictingEvidence []Evidence `json:"conflicting_evidence"`
	UnconfirmedEvidence []Evidence `json:"unconfirmed_evidence"`
}

// Explain is the endpoint that allows you to see how reported observations are linked with the final list of most probable conditions.
// For example, you can use this endpoint in the results page to display "reasons for" and "reasons against" particular conditions
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
