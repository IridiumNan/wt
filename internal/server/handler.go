package server

import (
	"log/slog"
	"net/http"

	"gitee.com/cai-zixiang_hainan/wt/internal/config"
	"gitee.com/cai-zixiang_hainan/wt/internal/model"
	"gitee.com/cai-zixiang_hainan/wt/internal/store"
	"gitee.com/cai-zixiang_hainan/wt/pkg/httphelper"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data := r.URL.Query()
	tokenHeadName := config.GetTokenHeadName(config.WTRead)
	readToken := r.Header.Get(tokenHeadName)

	// handle read token without access
	if !IsAccess(config.WTRead, readToken) {
		httphelper.SendJSONResponse(
			w,
			http.StatusForbidden,
			model.ForbiddenResponse("has no access to read"),
		)
		return
	}
	searchName := data.Get("name")

	slog.Debug("receive search patten (server)", "pattern", searchName)

	// handle empty search name
	if searchName == "" {
		httphelper.SendJSONResponse(
			w,
			http.StatusBadRequest,
			model.BadRequestResponse("package name required"),
		)
		return
	}

	results := store.SearchPackage(searchName)

	// handle not found
	if results == nil {
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("package not found"),
		)
		return
	}

	httphelper.SendJSONResponse(
		w,
		http.StatusOK,
		model.SuccessfulResponse(results, "package found"),
	)
}

func installHandler(w http.ResponseWriter, r *http.Request) {
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
}

func replaceHandler(w http.ResponseWriter, r *http.Request) {
}

func mvHandler(w http.ResponseWriter, r *http.Request) {
}

func rmHandler(w http.ResponseWriter, r *http.Request) {
}

func listHandler(w http.ResponseWriter, r *http.Request) {
}
