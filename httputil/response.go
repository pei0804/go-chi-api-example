package httputil

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, i interface{}) (err error) {
	b, err := json.Marshal(i)
	if err != nil {
		return
	}
	return JSONBlob(w, status, b)
}

var indent = "  "

func JSONPretty(w http.ResponseWriter, status int, i interface{}) (err error) {
	b, err := json.MarshalIndent(i, "", indent)
	if err != nil {
		return
	}
	return JSONBlob(w, status, b)
}

func JSONBlob(w http.ResponseWriter, status int, b []byte) (err error) {
	return Blob(w, status, MIMEApplicationJSONCharsetUTF8, b)
}

func String(w http.ResponseWriter, status int, s string) (err error) {
	return Blob(w, status, MIMETextPlainCharsetUTF8, []byte(s))
}

func Blob(w http.ResponseWriter, status int, contentType string, b []byte) (err error) {
	w.Header().Set(HeaderContentType, contentType)
	w.WriteHeader(status)
	_, err = w.Write(b)
	return
}
