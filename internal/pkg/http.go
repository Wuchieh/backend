package pkg

import (
	"github.com/goccy/go-json"
	"io"
	"net/http"
)

func DelCookie(w http.ResponseWriter, name string) {
	c := &http.Cookie{
		Name:   name,
		MaxAge: -1,
	}
	http.SetCookie(w, c)
}

func RespUnmarshal(h *http.Response, t interface{}) error {
	b, err := io.ReadAll(h.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, t)
}
