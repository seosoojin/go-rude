package rude

import "net/http"

func WriteErrorAsProblem(w http.ResponseWriter, r *http.Request, err error) {
	FromError(err).Write(w, r)
}

func WriteProblem(w http.ResponseWriter, r *http.Request, problem ProblemDetails) {
	problem.Write(w, r)
}
