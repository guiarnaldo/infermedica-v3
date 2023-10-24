package infermedica

import (
	"encoding/json"
	"net/http"
	"time"
)

type RiskFactorRes struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	CommonName          string    `json:"common_name"`
	Question            string    `json:"question"`
	QuestionThirdPerson string    `json:"question_third_person"`
	SexFilter           SexFilter `json:"sex_filter"`
	Category            string    `json:"category"`
	Extras              any       `json:"extras"`
	ImageURL            string    `json:"image_url"`
	ImageSource         string    `json:"image_source"`
}

func (a *App) RiskFactors(age int32) (*[]RiskFactorRes, error) {
	req, err := a.prepareRequest("GET", "risk_factors?age.value="+string(age), nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = checkResponse(res)

	if err != nil {
		return nil, err
	}
	r := []RiskFactorRes{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (a *App) RiskFactorByID(id string) (*RiskFactorRes, error) {
	req, err := a.prepareRequest("GET", "risk_factors/"+id, nil)
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
	r := RiskFactorRes{}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
