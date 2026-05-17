package server

import (
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
			model.ForbiddenResponse(errMsg),
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
	// Step 1: Get the boundary from Content-Type header
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/") {
		httphelper.SendJSONResponse(
			w,
			http.StatusBadRequest,
			model.BadRequestResponse("invalid content type"),
		)
		return
	}

	// Extract boundary (e.g., "boundary=----WebKitFormBoundary...")
	boundary := strings.TrimPrefix(contentType, "multipart/form-data; boundary=")
	if boundary == "" {
		httphelper.SendJSONResponse(
			w,
			http.StatusBadRequest,
			model.BadRequestResponse("missing boundary"),
		)
		return
	}

	// Step 2: Create a multipart reader with a small buffer (e.g., 64KB)
	reader := multipart.NewReader(r.Body, boundary)

	var pkgName string
	var filePart *multipart.Part

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			httphelper.SendJSONResponse(w, http.StatusBadRequest, model.BadRequestResponse("error reading form"))
		}

		if part.FormName() == "name" {

			nameBytes, _ := io.ReadAll(part)
			pkgName = string(nameBytes)
			slog.Debug("get the custom file name", "name", pkgName)
		} else if part.FormName() == "file" {
			filePart = part

			break
		}
		part.Close()
	}

	if filePart == nil {
		httphelper.SendJSONResponse(w, http.StatusBadRequest, model.BadRequestResponse("missing file filed"))
		return
	}

	defer filePart.Close()

	fileName := pkgName
	if fileName == "" {
		fileName = filePart.FormName()
	}

	savePath := filepath.Join(commonpresets.DataDir, fileName)
	dst, err := os.Create(savePath)
	if err != nil {
		httphelper.SendJSONResponse(w, http.StatusInternalServerError, model.InternalErrorResponse("fail to create file:"+savePath))
		return
	}

	defer dst.Close()

	written, err := io.Copy(dst, filePart)
	if err != nil {
		httphelper.SendJSONResponse(w, http.StatusInternalServerError, model.InternalErrorResponse("fail to write file:"+savePath))
	}

	// Step 6: Update metadata
	fileInfo, _ := dst.Stat()

	store.AddPackage(fileInfo)

	slog.Info("Package uploaded via stream", "name", fileName, "size", written)
	httphelper.SendJSONResponse(w, http.StatusOK, model.SuccessfulResponse("package uploaded successfully", ""))
}

func replaceHandler(w http.ResponseWriter, r *http.Request) {
}

func mvHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := r.URL.Query()
	tokenHeadName := config.GetTokenHeadName(model.WTWrite)

	errMsg := "has no access to write"
	auth := model.Auth{
		WtMethod: model.WTWrite,
		Token:    r.Header.Get(tokenHeadName),
		ErrMsg:   errMsg,
	}

	if !isAccess(auth) {
		httphelper.SendJSONResponse(
			w,
			http.StatusForbidden,
			model.ForbiddenResponse(errMsg),
		)
		return
	}

	oldName := data.Get("old_name")
	newName := data.Get("new_name")
	slog.Debug("receive old name and new name of package", "old_name", oldName, "new_name", newName)

	if oldName == "" || newName == "" {
		httphelper.SendJSONResponse(
			w,
			http.StatusBadRequest,
			model.BadRequestResponse("package name requeired"),
		)

		return
	}

	oldPath := filepath.Join(commonpresets.DataDir, oldName)
	newPath := filepath.Join(commonpresets.DataDir, newName)
	if err := os.Rename(oldPath, newPath); err != nil {
		if os.IsNotExist(err) {
			httphelper.SendJSONResponse(
				w,
				http.StatusNotFound,
				model.NotFoundResponse("package "+oldName+" not found"),
			)
			return
		}
		httphelper.SendJSONResponse(
			w,
			http.StatusInternalServerError,
			model.InternalErrorResponse("error: "+err.Error()),
		)
		return
	}

	store.RenamePackage(oldName, newName)

	httphelper.SendJSONResponse(
		w,
		http.StatusOK,
		model.SuccessfulResponse("mv the "+oldName+" to "+newName, ""),
	)
}

func rmHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := r.URL.Query()
	tokenHeadName := config.GetTokenHeadName(model.WTWrite)

	errMsg := "has no access to rm"
	auth := model.Auth{
		WtMethod: model.WTWrite,
		Token:    r.Header.Get(tokenHeadName),
		ErrMsg:   errMsg,
	}

	if !isAccess(auth) {
		httphelper.SendJSONResponse(
			w,
			http.StatusForbidden,
			model.ForbiddenResponse(errMsg),
		)
		return
	}

	pkgName := data.Get("name")
	slog.Debug("pkgName to rm ", "name", pkgName)

	if pkgName == "" {
		httphelper.SendJSONResponse(
			w,
			http.StatusBadRequest,
			model.BadRequestResponse("package name reqired"),
		)
		return
	}

	err := store.DeletePackageByName(pkgName)
	if err != nil {
		if os.IsNotExist(err) {
			httphelper.SendJSONResponse(
				w,
				http.StatusNotFound,
				model.NotFoundResponse("package "+pkgName+" not found"),
			)
			return
		}
		httphelper.SendJSONResponse(
			w,
			http.StatusInternalServerError,
			model.InternalErrorResponse("error: "+err.Error()),
		)
		return
	}

	httphelper.SendJSONResponse(
		w,
		http.StatusOK,
		model.SuccessfulResponse("rm the "+pkgName+" success", ""),
	)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query()
	tokenHeadName := config.GetTokenHeadName(model.WTRead)

	errMsg := "has no access to read"
	auth := model.Auth{
		WtMethod: model.WTRead,
		Token:    r.Header.Get(tokenHeadName),
		ErrMsg:   errMsg,
	}

	if !isAccess(auth) {
		httphelper.SendJSONResponse(
			w,
			http.StatusForbidden,
			model.ForbiddenResponse(errMsg),
		)
		return
	}

	targetTag := data.Get("tag")

	slog.Debug("receive tag (server)", "tag", targetTag)

	errMsg = "require tag"
	if targetTag == "" {
		httphelper.SendJSONResponse(
			w,
			http.StatusBadRequest,
			model.BadRequestResponse(errMsg),
		)
	}

	nameList := store.ListPackagesByTag(targetTag)

	if nameList == nil {
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("there is no package tag as "+targetTag),
		)
		return
	}

	httphelper.SendJSONResponse(
		w,
		http.StatusOK,
		model.SuccessfulResponse(nameList, "packages found"),
	)
}

func syncHandler(w http.ResponseWriter, r *http.Request) {
	tokenHeadName := config.GetTokenHeadName(model.WTWrite)

	errMsg := "has no access to write"
	auth := model.Auth{
		WtMethod: model.WTWrite,
		Token:    r.Header.Get(tokenHeadName),
		ErrMsg:   errMsg,
	}

	if !isAccess(auth) {
		httphelper.SendJSONResponse(
			w,
			http.StatusForbidden,
			model.ForbiddenResponse(errMsg),
		)
		slog.Debug("check over")
		return
	}

	slog.Debug("sync the meta data")
	store.SyncMetaDataFromDisk()
	httphelper.SendJSONResponse(
		w,
		http.StatusOK,
		model.SuccessfulResponse("sync data from disk success", ""),
	)
}
