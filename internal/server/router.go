package server

import (
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := *http.NewServeMux()

	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/info", infoHandler)
	mux.HandleFunc("/install", installHandler)
	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/mv", mvHandler)
	mux.HandleFunc("/rm", rmHandler)
	mux.HandleFunc("/list", listHandler)
	mux.HandleFunc("/sync", syncHandler)
	mux.HandleFunc("/tag/list", tagListHandler)
	// mux.HandleFunc("/tag/add", addTagHandler)

	return &mux
}
