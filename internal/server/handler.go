package server

import (
	"log/slog"
	"net/http"
	"path/filepath"

	"gitee.com/cai-zixiang_hainan/wt/internal/config"
	"gitee.com/cai-zixiang_hainan/wt/internal/model"
	"gitee.com/cai-zixiang_hainan/wt/internal/presets/commonpresets"
	"gitee.com/cai-zixiang_hainan/wt/internal/store"
	"gitee.com/cai-zixiang_hainan/wt/pkg/httphelper"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data := r.URL.Query()
	tokenHeadName := config.GetTokenHeadName(model.WTRead)

	errMsg := "has no access to read"
	auth := model.Auth{
		WtMethod: model.WTRead,
		Token:    r.Header.Get(tokenHeadName),
		ErrMsg:   errMsg,
	}

	// handle read token without access
	if !isAccess(auth) {
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

func infoHandler(w http.ResponseWriter, r *http.Request) {
	searchHandler(w, r)
}

func installHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data := r.URL.Query()
	tokenHeadName := config.GetTokenHeadName(model.WTInstall)
	installToken := r.Header.Get(tokenHeadName)
	errMsg := "has no access to install"

	auth := model.Auth{
		WtMethod: model.WTInstall,
		Token:    installToken,
		ErrMsg:   errMsg,
	}

	if !isAccess(auth) {
		httphelper.SendJSONResponse(
			w,
			http.StatusForbidden,
			model.ForbiddenResponse("has no access to install"),
		)
		return
	}

	pkgName := data.Get("name")
	if filepath.Base(pkgName) != pkgName {
		httphelper.SendJSONResponse(w, http.StatusBadRequest, model.BadRequestResponse("invalid package name"))
		return
	}
	_, err := store.GetPackageInfo(pkgName)
	if err != nil {
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("pkg "+pkgName+" not found"),
		)
		return
	}

	filePath := filepath.Join(commonpresets.DataDir, pkgName)

	http.ServeFile(w, r, filePath)
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
