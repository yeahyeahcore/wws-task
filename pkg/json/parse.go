package json

import (
	"encoding/json"
	"io"
)

func Parse[T interface{}](reader io.Reader) (*T, error) {
	var jsonData T

	if err := json.NewDecoder(reader).Decode(&jsonData); err != nil {
		return nil, err
	}

	return &jsonData, nil
}
