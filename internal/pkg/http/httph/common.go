package httph

import (
	"encoding/json"
	"net/http"
)

func SendRaw(w http.ResponseWriter, statusCode int, mimeType string, data []byte) {
	if mimeType != "" {
		w.Header().Set(HeaderContentType, mimeType)
	}
	w.WriteHeader(statusCode)
	if len(data) > 0 {
		_, _ = w.Write(data)
	}
}

func SendEmpty(w http.ResponseWriter, statusCode int) {
	SendRaw(w, statusCode, "", nil)
}

func SendEncodedWithMIME(w http.ResponseWriter, r *http.Request, statusCode int, mimetype string, obj any) {
	data, err := json.Marshal(obj)
	if err != nil {
		ErrorApply(w, http.StatusInternalServerError, err.Error())
		return
	}
	SendRaw(w, statusCode, mimetype, data)
}

func SendEncoded(w http.ResponseWriter, r *http.Request, statusCode int, obj any) {
	SendEncodedWithMIME(w, r, statusCode, MIMEApplicationJSONCharsetUTF8, obj)
}
