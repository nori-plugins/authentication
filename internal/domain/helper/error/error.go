package error

import "net/http"

type ErrorHelper interface {
	Error(w http.ResponseWriter, err error)
}