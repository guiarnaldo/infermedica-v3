package infermedica

import (
	"encoding/json"
	"net/http"
	"time"
)

type ParseReq struct {
	Text string `json:"text"`
	Age  Age    `json:"age"`
	Sex  Sex    `json:"sex,omitempty"`
}

type ParseRes struct {
	Mentions []struct {
		ID         string `json:"id"`
		Orth       string `json:"orth"`
		ChoiceID   string `json:"choice_id"`
		Name       string `json:"name"`
		CommonName string `json:"common_name"`
		Type       string `json:"type"`
	} `json:"mentions"`
	Obvious bool `json:"obvious"`
}

func (a *App) Parse(pr ParseReq) (*ParseRes, error) {
	// Required to use "infermedica-en" model, because NPL is only avaliable in english at the moment
	model := a.model
	a.model = ""

	req, err := a.preparePOSTRequest("parse", pr)
	if err != nil {
		return nil, err
	}
	a.model = model
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
	r := ParseRes{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// Converts a Parse Response into an Evidence
func ParseToEvidence(p *ParseRes) (evidences []Evidence) {
	var e Evidence

	for i := range p.Mentions {
		e.ChoiceID = EvidenceChoiceID(p.Mentions[i].ChoiceID)
		e.ID = p.Mentions[i].ID
		evidences = append(evidences, e)
	}
	return
}
