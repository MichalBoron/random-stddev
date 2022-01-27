package main

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
