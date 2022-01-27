package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

//ErrorStruct represents error details returned in the body of an http response
type ErrorStruct struct {
	Timestamp string `json:"timestamp,omitempty"`
	Status    int    `json:"status,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
	Message   string `json:"message,omitempty"`
	Path      string `json:"path,omitempty"`
}

func makeErrorStruct(status int, errorMsg, message, path string) ErrorStruct {
	return ErrorStruct{
		Timestamp: Now().String(),
		Status:    status,
		ErrorMsg:  errorMsg,
		Message:   message,
		Path:      path,
	}
}

//writeAsJsonWithStatus writes JSON representation of v to w
// and writes header with given status.
//If conversion to JSON fails, http.InternalServerError status is set
// and no body is sent.
func writeAsJsonWithStatus(v interface{}, status int, w http.ResponseWriter) {
	body, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(body)
}

func areParamsPresent(params url.Values, cases ...string) bool {
	for _, c := range cases {
		_, isPresent := params[c]
		if !isPresent {
			return false
		}
	}
	return true
}
