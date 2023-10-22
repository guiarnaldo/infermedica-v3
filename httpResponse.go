package infermedica

import (
	"fmt"
	"net/http"
)

// Check response status, return a error with the status code if not it is not 200
func checkResponse(res *http.Response) error{
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", res.Status)
	}
	return nil
}