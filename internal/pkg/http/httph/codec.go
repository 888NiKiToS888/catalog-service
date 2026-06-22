package httph

import (
	"encoding/json"
	"io"
	"net/http"
)

func EncodeJSON(w io.Writer, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}

func DecodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
