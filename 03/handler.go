package main

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
)

var logger = MustGetLogger("api")

type handler func(http.ResponseWriter, *http.Request) (int, interface{}, error)

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rv := recover(); rv != nil {
			log.Print(rv)
			debug.PrintStack()
			logger.Errorf("panic: %s", rv)
			http.Error(w, http.StatusText(
				http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}()
	status, res, err := h(w, r)
	if err != nil {
		logger.Infof("error: %s", err)
		respondError(w, status, err)
		return
	}
	respondJSON(w, status, res)
	return
}

// respondJSON レスポンスとして返すjsonを生成して、writerに書き込む
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "    ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// respondError レスポンスとして返すエラーを生成する
func respondError(w http.ResponseWriter, code int, err error) {
	log.Printf("err: %v", err)
	if e, ok := err.(*HTTPError); ok {
		respondJSON(w, e.Code, e)
	} else if err != nil {
		he := HTTPError{
			Code:    code,
			Message: err.Error(),
		}
		respondJSON(w, code, he)
	}
}
