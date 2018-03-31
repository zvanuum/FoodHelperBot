package util

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

func UnmarshalBody(body io.ReadCloser, target interface{}) error {
	var bodyBytes []byte
	var err error

	if bodyBytes, err = ioutil.ReadAll(body); err == nil {
		if err2 := json.Unmarshal(bodyBytes, &target); err2 != nil {
			return fmt.Errorf("failed unmarshalling body: %s", err2)
		}

		return nil
	}

	return fmt.Errorf("failed to read body: %s", err)
}
