package request

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

// GetJSON request Get JSON
func GetJSON(url string, target interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// SendJSON request Send JSON
func SendJSON(url string, data interface{}) error {
	json, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewReader(json))
	req.Header.Set("Content-Type", "application/json")

	r, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}
