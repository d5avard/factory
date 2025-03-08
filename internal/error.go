package internal

import "net/http"

func HttpError(w http.ResponseWriter, r *http.Request, message string, code int) {
	LogRequest(r, message, code)
	http.Error(w, message, code)
}
