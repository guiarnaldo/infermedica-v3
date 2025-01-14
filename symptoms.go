package infermedica

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type SymptomRes struct {
	ID                  string         `json:"id"`
	Name                string         `json:"name"`
	CommonName          string         `json:"common_name"`
	Category            string         `json:"category"`
	Seriousness         string         `json:"seriousness"`
	Children            []SymptomChild `json:"children"`
	ImageURL            string         `json:"image_url"`
	ImageSource         string         `json:"image_source"`
	ParentID            string         `json:"parent_id"`
	ParentRelation      string         `json:"parent_relation"`
	Question            string         `json:"question"`
	QuestionThirdPerson string         `json:"question_third_person"`
	SexFilter           SexFilter      `json:"sex_filter"`
	Extra               any            `json:"extra"`
}

type SymptomChild struct {
	ID             string `json:"id"`
	ParentRelation string `json:"parent_relation"`
}

func (a *App) Symptoms(age Age, enableTriage3 bool) (*[]SymptomRes, error) {
	req, err := a.prepareRequest("GET", "symptoms?age.value="+strconv.Itoa(age.Value)+"&age.unit"+string(age.Unit)+"&enableTriage3="+strconv.FormatBool(enableTriage3), nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 10,
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

	var r []SymptomRes
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (a *App) SymptomByID(id string, age Age, enableTriage3 bool) (*SymptomRes, error) {
	req, err := a.prepareRequest("GET", "symptoms/"+id+"?age.value="+strconv.Itoa(age.Value)+"&age.unit"+string(age.Unit)+"&enableTriage3="+strconv.FormatBool(enableTriage3), nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
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

	var r SymptomRes
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
