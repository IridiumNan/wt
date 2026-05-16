package server

import "net/http"

func NewRouter() *http.ServeMux {
	mux := *http.NewServeMux()

	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/info", searchHandler)
	mux.HandleFunc("/install", installHandler)
	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/mv", mvHandler)
	mux.HandleFunc("/rm", rmHandler)

	return &mux
}
