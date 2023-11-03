package infermedica

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

// Check response status, returns an error with the status code if not it is not 200
func checkResponse(res *http.Response) error {
	if res.StatusCode != http.StatusOK {
		var response Response
		json.NewDecoder(res.Body).Decode(&response)
		return fmt.Errorf("%s: %s", res.Status, response.Message)
	}
	return nil
}
