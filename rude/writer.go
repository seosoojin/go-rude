package rude

import "net/http"

func WriteProblem(w http.ResponseWriter, r *http.Request, err error) {
	FromError(err).Write(w, r)
}
