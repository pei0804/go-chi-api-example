package httputil

import "fmt"

type HTTPError struct {
	Status  int
	Message string
}

func (he *HTTPError) Error() string {
	return fmt.Sprintf("status=%d, message=%v", he.Status, he.Message)
}

func NewHTTPError(status int, message ...interface{}) *HTTPError {
	he := &HTTPError{Status: status}
	return he
}