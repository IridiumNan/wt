package server

import (
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gitee.com/cai-zixiang_hainan/wt/internal/config"
	"gitee.com/cai-zixiang_hainan/wt/internal/model"
	"gitee.com/cai-zixiang_hainan/wt/internal/store"
	"gitee.com/cai-zixiang_hainan/wt/pkg/httphelper"
)

// searchHandler : search package in available scope, tag is not required
func searchHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := r.URL.Query()
	searchName := data.Get("name")

	// if !TokenOk(w, r, model.WTRead, "") {
	// 	return
	// }

	clientToken := r.Header.Get(config.GetTokenHeadName(model.WTRead))
	// search in the available scope
	tags := GetAvailableTags(clientToken, model.WTRead)
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

	allResults := store.SearchPackage(searchName)
	var results []*model.Package
	for i := range allResults {
		if slices.Contains(tags, allResults[i].Tag) {
			results = append(results, allResults[i])
		}
	}

	slog.Debug("check search results", "available tags", tags, "allresults", allResults, "results", results)
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

// infoHandler : get information of package, and the package tag will handle by server so tag is no required
func infoHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := r.URL.Query()
	pkgName := data.Get("name")

	pkg, err := store.GetPackage(pkgName)
	if err == nil && !TokenOk(w, r, model.WTRead, pkg.Tag) {
		return
	}
	if err != nil {
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("package "+pkgName+" not found"),
		)
	}

	httphelper.SendJSONResponse(
		w,
		http.StatusOK,
		model.SuccessfulResponse([]*model.Package{pkg}, "package found"),
	)
}

// installHandler : handle the install reqest and the tag is not required
func installHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data := r.URL.Query()
	pkgName := data.Get("name")

	pkg, err := store.GetPackage(pkgName)
	if err != nil {
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("pkg "+pkgName+" not found"),
		)
		return
	}
	if filepath.Base(pkgName) != pkgName {
		httphelper.SendJSONResponse(w, http.StatusBadRequest, model.BadRequestResponse("invalid package name"))
		return
	}
	tag := pkg.Tag

	if !TokenOk(w, r, model.WTInstall, tag) {
		return
	}

	filePath := filepath.Join(config.DataDir, pkgName)

	http.ServeFile(w, r, filePath)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if !TokenOk(w, r, model.WTWrite, config.DefaultTagTemp) {
		return
	}
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

	savePath := filepath.Join(config.DataDir, fileName)

	backupPath := ""
	if _, err := os.Stat(savePath); err == nil {
		slog.Info("Package exist try to replace", "packageName", fileName)

		backupPath = savePath + ".bak"
		if err = os.Rename(savePath, backupPath); err != nil {
			slog.Warn("fail to rename old package to backup package", "error", err)
		}
	}

	dst, err := os.Create(savePath)
	if err != nil {
		if _, err = os.Stat(backupPath); err == nil {
			err = os.Rename(backupPath, savePath)
			if err != nil {
				slog.Info("recover the old package successfully", "pkgName", fileName)
			}
		}
		httphelper.SendJSONResponse(w, http.StatusInternalServerError, model.InternalErrorResponse("fail to create file:"+savePath))
		return
	}

	defer dst.Close()

	written, err := io.Copy(dst, filePart)
	if err != nil {
		dst.Close()
		os.Remove(savePath)

		if backupPath != "" {
			if rollbakErr := os.Rename(backupPath, savePath); rollbakErr == nil {
				slog.Info("recover old package after write failure", "pkgName", fileName)
			}
		}
		httphelper.SendJSONResponse(w, http.StatusInternalServerError, model.InternalErrorResponse("fail to write file:"+savePath))
		return
	}

	if _, err = os.Stat(backupPath); err == nil {
		slog.Info("found backup file, try to remove", "backupName", backupPath)
		err = os.Remove(backupPath)
		if err != nil {
			slog.Error("fail to remove backup file", "err", err.Error())
		}
	}

	// Step 6: Update metadata
	fileInfo, _ := dst.Stat()

	store.AddPackage(fileInfo)

	slog.Info("Package uploaded via stream", "name", fileName, "size", written)
	httphelper.SendJSONResponse(w, http.StatusOK, model.SuccessfulResponse("package uploaded successfully", ""))
}

// mvHandler : rename the package and the tag is not required from header
func mvHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := r.URL.Query()

	oldName := data.Get("old_name")
	newName := data.Get("new_name")

	if oldName == "" || newName == "" {
		httphelper.SendJSONResponse(
			w,
			http.StatusBadRequest,
			model.BadRequestResponse("package name requeired"),
		)

		return
	}

	oldPkg, err := store.GetPackage(oldName)
	if err != nil {
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("package "+oldName+" not found"),
		)
		return
	}

	if !TokenOk(w, r, model.WTWrite, oldPkg.Tag) {
		return
	}

	slog.Debug("receive old name and new name of package", "old_name", oldName, "new_name", newName)

	oldPath := filepath.Join(config.DataDir, oldName)
	newPath := filepath.Join(config.DataDir, newName)
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

// rmHandler : handle the remove request and the tag is not reqired from the header
func rmHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := r.URL.Query()
	pkgName := data.Get("name")

	pkg, err := store.GetPackage(pkgName)
	if err != nil {
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("package "+pkgName+" not found"),
		)
		return
	}
	if !TokenOk(w, r, model.WTWrite, pkg.Tag) {
		return
	}

	slog.Debug("pkgName to rm ", "name", pkgName)

	if pkgName == "" {
		httphelper.SendJSONResponse(
			w,
			http.StatusBadRequest,
			model.BadRequestResponse("package name reqired"),
		)
		return
	}

	err = store.DeletePackageByName(pkgName)
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

// listHandler : handle the list require and the tag is required
func listHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := r.URL.Query()
	tag := data.Get("tag")
	// check the token -> global token -> tag token
	if !TokenOk(w, r, model.WTRead, tag) {
		return
	}

	slog.Debug("receive tag (server)", "tag", tag)

	errMsg := "require tag"
	if tag == "" {
		httphelper.SendJSONResponse(
			w,
			http.StatusBadRequest,
			model.BadRequestResponse(errMsg),
		)
	}

	nameList := store.ListPackagesByTag(tag)

	if nameList == nil {
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("there is no package tag as "+tag),
		)
		return
	}

	httphelper.SendJSONResponse(
		w,
		http.StatusOK,
		model.SuccessfulResponse(nameList, "packages found"),
	)
}

// syncHandler : handle the sync request which the token has access to write and tag is not required for header
func syncHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if !TokenOk(w, r, model.WTWrite, "") {
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

// tagListHandler : list all tags the which the client has access to read
func tagListHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	clientToken := r.Header.Get(config.GetTokenHeadName(model.WTRead))
	tags := GetAvailableTags(clientToken, model.WTRead)

	slog.Debug("get tag list", "client_token", clientToken, "tags", tags)

	msg := fmt.Sprintf(
		"has access to %d tags as below",
		len(tags),
	)
	httphelper.SendJSONResponse(
		w,
		http.StatusOK,
		model.SuccessfulResponse(tags, msg),
	)
}

func addTagHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := r.URL.Query()
	tagName := data.Get("tag")

	if !TokenOk(w, r, model.WTWrite, "") {
		return
	}

	config.AddTagTokenList(tagName)

	slog.Info("add a new tag", "tagName", tagName)

	httphelper.SendJSONResponse(
		w,
		http.StatusOK,
		model.SuccessfulResponse("add the tag:"+tagName+"successfully", ""),
	)
}

func tagUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer r.Body.Close()
	data := r.URL.Query()
	newTag := data.Get("new_tag")
	pkgName := data.Get("name")

	slog.Debug("update request check", "pkg_name", pkgName, "new_tag", newTag)

	var pkg *model.Package
	pkg, err = store.GetPackage(pkgName)

	slog.Debug("check pkg", "pkgInfo", pkg.Name)
	if err != nil {
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("package "+pkgName+" not found"),
		)
		return
	}
	oldTag := pkg.Tag

	if !TokenOk(w, r, model.WTWrite, oldTag) || !TokenOk(w, r, model.WTWrite, newTag) {
		return
	}

	err = store.UpdateTag(pkgName, newTag)
	if err != nil {
		httphelper.SendJSONResponse(
			w,
			http.StatusInternalServerError,
			model.InternalErrorResponse("fail to update the tag:"+err.Error()),
		)
		return
	}

	httphelper.SendJSONResponse(
		w,
		http.StatusOK,
		model.SuccessfulResponse("update the package "+pkgName+" from old tag "+oldTag+" to new tag "+newTag+" successfully", ""),
	)
}

func tagRmHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer r.Body.Close()
	data := r.URL.Query()
	tag := data.Get("tag")

	slog.Debug("tag rm request check", "tag", tag)

	if !TokenOk(w, r, model.WTWrite, "") {
		return
	}
	tagList := config.GetTagList()

	if !slices.Contains(tagList, tag) {
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("tag "+tag+" not found"),
		)
		return

	}

	err = store.RemoveTag(tag)
	if err != nil {
		httphelper.SendJSONResponse(
			w,
			http.StatusInternalServerError,
			model.InternalErrorResponse("error when remove the tag: "+err.Error()),
		)
	}
}
