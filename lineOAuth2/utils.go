package lineOAuth2

import (
	"encoding/json"
	"io"
	"net/http"
)

func RespUnmarshal(h *http.Response, t interface{}) error {
	b, err := io.ReadAll(h.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, t)
}
