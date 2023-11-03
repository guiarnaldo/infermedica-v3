package infermedica

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LookupRes struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

func (a *App) Lookup(phrase string, sex Sex) (*LookupRes, error) {
	if sex.IsValid() != nil {
		return nil, fmt.Errorf("infermedica: Unexpected value for Sex")
	}
	url := "lookup?phrase=" + phrase + "&sex=" + sex.String()
	req, err := a.prepareRequest("GET", url, nil)
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
	var r LookupRes
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
