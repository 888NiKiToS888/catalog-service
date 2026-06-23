package httph

import "net/http"

type Error struct {
	Message string `json:"error"`
}

func ErrorApply(w http.ResponseWriter, code int, message string) {
	w.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	w.WriteHeader(code)
	_ = EncodeJSON(w, Error{message})
}
