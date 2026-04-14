package rude

import "net/http"

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	From(err).Write(w, r)
}
