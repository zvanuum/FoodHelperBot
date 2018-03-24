package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func UnmarshalResponse(res *http.Response, target interface{}) error {
	var body []byte
	var err error

	if body, err = ioutil.ReadAll(res.Body); err == nil {
		if err2 := json.Unmarshal(body, &target); err2 != nil {
			return fmt.Errorf("failed unmarshalling response: %s", err2)
		}

		return nil
	}

	return fmt.Errorf("failed to read response body: %s", err)
}
