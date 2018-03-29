package util

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

func UnmarshalResponse(resBody io.ReadCloser, target interface{}) error {
	var body []byte
	var err error

	if body, err = ioutil.ReadAll(resBody); err == nil {
		if err2 := json.Unmarshal(body, &target); err2 != nil {
			return fmt.Errorf("failed unmarshalling response: %s", err2)
		}

		return nil
	}

	return fmt.Errorf("failed to read response body: %s", err)
}
