package infermedica

import (
	"encoding/json"
	"net/http"
	"time"
)

type ConceptsRes struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	CommonName string `json:"common_name"`
}

// Concepts returns all concepts
func (a *App) Concepts() (*[]ConceptsRes, error) {
	req, err := a.prepareRequest("GET", "concepts", nil)
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
	var r []ConceptsRes
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (a *App) ConceptsByID(id string) (*ConceptsRes, error) {
	req, err := a.prepareRequest("GET", "concepts/"+id, nil)
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
	var r ConceptsRes
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
