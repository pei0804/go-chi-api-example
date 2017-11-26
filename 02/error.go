package main

import "fmt"

// HTTPError エラー用
type HTTPError struct {
	Code    int
	Message string
}

func (he *HTTPError) Error() string {
	return fmt.Sprintf("code=%d, message=%v", he.Code, he.Message)
}
