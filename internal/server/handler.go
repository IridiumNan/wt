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

	slog.Info("Search request received", "method", "GET", "path", "/search", "query_name", searchName, "remote_addr", r.RemoteAddr)

	// if !TokenOk(w, r, model.WTRead, "") {
	// 	return
	// }

	clientToken := r.Header.Get(config.GetTokenHeadName(model.WTRead))
	// search in the available scope
	tags := GetAvailableTags(clientToken, model.WTRead)
	slog.Debug("Search token validation", "token_present", clientToken != "", "available_tags", tags)

	// handle empty search name
	if searchName == "" {
		slog.Warn("Search failed: empty package name", "remote_addr", r.RemoteAddr)
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

	slog.Debug("Search results filtered", "total_found", len(allResults), "filtered_count", len(results), "available_tags", tags)
	// handle not found
	if results == nil {
		slog.Info("Search returned no results", "pattern", searchName, "remote_addr", r.RemoteAddr)
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("package not found"),
		)
		return
	}

	slog.Info("Search completed successfully", "pattern", searchName, "results_count", len(results), "remote_addr", r.RemoteAddr)
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

	slog.Info("Info request received", "method", "GET", "path", "/info", "package_name", pkgName, "remote_addr", r.RemoteAddr)

	pkg, err := store.GetPackage(pkgName)
	if err == nil && !TokenOk(w, r, model.WTRead, pkg.Tag) {
		slog.Warn("Info request denied: insufficient read permission", "package_name", pkgName, "tag", pkg.Tag, "remote_addr", r.RemoteAddr)
		return
	}
	if err != nil {
		slog.Info("Info request failed: package not found", "package_name", pkgName, "remote_addr", r.RemoteAddr)
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("package "+pkgName+" not found"),
		)
		return
	}

	slog.Info("Info request completed", "package_name", pkgName, "tag", pkg.Tag, "size", pkg.Size, "remote_addr", r.RemoteAddr)
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

	slog.Info("Install request received", "method", "GET", "path", "/install", "package_name", pkgName, "remote_addr", r.RemoteAddr)

	pkg, err := store.GetPackage(pkgName)
	if err != nil {
		slog.Info("Install request failed: package not found", "package_name", pkgName, "remote_addr", r.RemoteAddr)
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("pkg "+pkgName+" not found"),
		)
		return
	}
	if filepath.Base(pkgName) != pkgName {
		slog.Warn("Install request failed: invalid package name", "package_name", pkgName, "remote_addr", r.RemoteAddr)
		httphelper.SendJSONResponse(w, http.StatusBadRequest, model.BadRequestResponse("invalid package name"))
		return
	}
	tag := pkg.Tag

	if !TokenOk(w, r, model.WTInstall, tag) {
		slog.Warn("Install request denied: insufficient install permission", "package_name", pkgName, "tag", tag, "remote_addr", r.RemoteAddr)
		return
	}

	filePath := filepath.Join(config.DataDir, pkgName)

	slog.Info("Serving package file", "package_name", pkgName, "file_path", filePath, "tag", tag, "remote_addr", r.RemoteAddr)
	http.ServeFile(w, r, filePath)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Upload request received", "method", "POST", "path", "/upload", "remote_addr", r.RemoteAddr)

	if !TokenOk(w, r, model.WTWrite, config.DefaultTagTemp) {
		slog.Warn("Upload request denied: insufficient write permission", "remote_addr", r.RemoteAddr)
		return
	}
	// Step 1: Get the boundary from Content-Type header
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/") {
		slog.Warn("Upload failed: invalid content type", "content_type", contentType, "remote_addr", r.RemoteAddr)
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
		slog.Warn("Upload failed: missing boundary", "remote_addr", r.RemoteAddr)
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
			slog.Error("Upload failed: error reading form", "error", err.Error(), "remote_addr", r.RemoteAddr)
			httphelper.SendJSONResponse(w, http.StatusBadRequest, model.BadRequestResponse("error reading form"))
			return
		}

		if part.FormName() == "name" {

			nameBytes, _ := io.ReadAll(part)
			pkgName = string(nameBytes)
			slog.Debug("Upload: custom filename provided", "custom_name", pkgName)
		} else if part.FormName() == "file" {
			filePart = part

			break
		}
		part.Close()
	}

	if filePart == nil {
		slog.Warn("Upload failed: missing file field", "remote_addr", r.RemoteAddr)
		httphelper.SendJSONResponse(w, http.StatusBadRequest, model.BadRequestResponse("missing file filed"))
		return
	}

	defer filePart.Close()

	fileName := pkgName
	if fileName == "" {
		fileName = filePart.FileName()
		slog.Debug("Upload: using original filename", "original_name", fileName)
	}

	savePath := filepath.Join(config.DataDir, fileName)

	backupPath := ""
	if _, err := os.Stat(savePath); err == nil {
		slog.Info("Package exists, will replace", "package_name", fileName, "existing_path", savePath)

		backupPath = savePath + ".bak"
		if err = os.Rename(savePath, backupPath); err != nil {
			slog.Warn("Failed to backup existing package", "error", err.Error(), "package_name", fileName)
		}
	}

	dst, err := os.Create(savePath)
	if err != nil {
		slog.Error("Upload failed: cannot create file", "error", err.Error(), "save_path", savePath)
		if _, err = os.Stat(backupPath); err == nil {
			err = os.Rename(backupPath, savePath)
			if err != nil {
				slog.Info("Recovered old package successfully", "pkg_name", fileName)
			}
		}
		httphelper.SendJSONResponse(w, http.StatusInternalServerError, model.InternalErrorResponse("fail to create file:"+savePath))
		return
	}

	defer dst.Close()

	written, err := io.Copy(dst, filePart)
	if err != nil {
		slog.Error("Upload failed: write error", "error", err.Error(), "bytes_written", written, "package_name", fileName)
		dst.Close()
		os.Remove(savePath)

		if backupPath != "" {
			if rollbakErr := os.Rename(backupPath, savePath); rollbakErr == nil {
				slog.Info("Recovered old package after write failure", "pkg_name", fileName)
			}
		}
		httphelper.SendJSONResponse(w, http.StatusInternalServerError, model.InternalErrorResponse("fail to write file:"+savePath))
		return
	}

	if _, err = os.Stat(backupPath); err == nil {
		slog.Info("Removing backup file", "backup_path", backupPath)
		err = os.Remove(backupPath)
		if err != nil {
			slog.Error("Failed to remove backup file", "error", err.Error(), "backup_path", backupPath)
		}
	}

	// Step 6: Update metadata
	fileInfo, _ := dst.Stat()

	store.AddPackage(fileInfo)

	slog.Info("Package uploaded successfully", "name", fileName, "size_bytes", written, "size_human", fmt.Sprintf("%.2f MB", float64(written)/1024/1024), "tag", config.DefaultTagTemp, "path", savePath)
	httphelper.SendJSONResponse(w, http.StatusOK, model.SuccessfulResponse("package uploaded successfully", ""))
}

// mvHandler : rename the package and the tag is not required from header
func mvHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := r.URL.Query()

	oldName := data.Get("old_name")
	newName := data.Get("new_name")

	slog.Info("Rename request received", "method", "POST", "path", "/mv", "old_name", oldName, "new_name", newName, "remote_addr", r.RemoteAddr)

	if oldName == "" || newName == "" {
		slog.Warn("Rename failed: missing package names", "old_name", oldName, "new_name", newName, "remote_addr", r.RemoteAddr)
		httphelper.SendJSONResponse(
			w,
			http.StatusBadRequest,
			model.BadRequestResponse("package name requeired"),
		)

		return
	}

	oldPkg, err := store.GetPackage(oldName)
	if err != nil {
		slog.Info("Rename failed: old package not found", "old_name", oldName, "remote_addr", r.RemoteAddr)
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("package "+oldName+" not found"),
		)
		return
	}

	if !TokenOk(w, r, model.WTWrite, oldPkg.Tag) {
		slog.Warn("Rename denied: insufficient write permission", "old_name", oldName, "tag", oldPkg.Tag, "remote_addr", r.RemoteAddr)
		return
	}

	slog.Debug("Renaming package files", "old_path", filepath.Join(config.DataDir, oldName), "new_path", filepath.Join(config.DataDir, newName))

	oldPath := filepath.Join(config.DataDir, oldName)
	newPath := filepath.Join(config.DataDir, newName)
	if err := os.Rename(oldPath, newPath); err != nil {
		if os.IsNotExist(err) {
			slog.Info("Rename failed: file not found on disk", "old_name", oldName, "remote_addr", r.RemoteAddr)
			httphelper.SendJSONResponse(
				w,
				http.StatusNotFound,
				model.NotFoundResponse("package "+oldName+" not found"),
			)
			return
		}
		slog.Error("Rename failed: filesystem error", "error", err.Error(), "old_name", oldName, "new_name", newName)
		httphelper.SendJSONResponse(
			w,
			http.StatusInternalServerError,
			model.InternalErrorResponse("error: "+err.Error()),
		)
		return
	}

	store.RenamePackage(oldName, newName)

	slog.Info("Package renamed successfully", "old_name", oldName, "new_name", newName, "tag", oldPkg.Tag)
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

	slog.Info("Delete request received", "method", "DELETE", "path", "/rm", "package_name", pkgName, "remote_addr", r.RemoteAddr)

	pkg, err := store.GetPackage(pkgName)
	if err != nil {
		slog.Info("Delete failed: package not found", "package_name", pkgName, "remote_addr", r.RemoteAddr)
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("package "+pkgName+" not found"),
		)
		return
	}
	if !TokenOk(w, r, model.WTWrite, pkg.Tag) {
		slog.Warn("Delete denied: insufficient write permission", "package_name", pkgName, "tag", pkg.Tag, "remote_addr", r.RemoteAddr)
		return
	}

	slog.Debug("Deleting package from store", "package_name", pkgName, "tag", pkg.Tag)

	if pkgName == "" {
		slog.Warn("Delete failed: empty package name", "remote_addr", r.RemoteAddr)
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
			slog.Info("Delete failed: file not found on disk", "package_name", pkgName, "remote_addr", r.RemoteAddr)
			httphelper.SendJSONResponse(
				w,
				http.StatusNotFound,
				model.NotFoundResponse("package "+pkgName+" not found"),
			)
			return
		}
		slog.Error("Delete failed: filesystem error", "error", err.Error(), "package_name", pkgName)
		httphelper.SendJSONResponse(
			w,
			http.StatusInternalServerError,
			model.InternalErrorResponse("error: "+err.Error()),
		)
		return
	}

	slog.Info("Package deleted successfully", "package_name", pkgName, "tag", pkg.Tag)
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

	slog.Info("List request received", "method", "GET", "path", "/list", "tag", tag, "remote_addr", r.RemoteAddr)

	// check the token -> global token -> tag token
	if !TokenOk(w, r, model.WTRead, tag) {
		slog.Warn("List denied: insufficient read permission", "tag", tag, "remote_addr", r.RemoteAddr)
		return
	}

	slog.Debug("Listing packages for tag", "tag", tag)

	errMsg := "require tag"
	if tag == "" {
		slog.Warn("List failed: tag parameter required", "remote_addr", r.RemoteAddr)
		httphelper.SendJSONResponse(
			w,
			http.StatusBadRequest,
			model.BadRequestResponse(errMsg),
		)
		return
	}

	nameList := store.ListPackagesByTag(tag)

	if nameList == nil {
		slog.Info("List returned no packages", "tag", tag, "remote_addr", r.RemoteAddr)
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("there is no package tag as "+tag),
		)
		return
	}

	slog.Info("List completed successfully", "tag", tag, "package_count", len(nameList))
	httphelper.SendJSONResponse(
		w,
		http.StatusOK,
		model.SuccessfulResponse(nameList, "packages found"),
	)
}

// syncHandler : handle the sync request which the token has access to write and tag is not required for header
func syncHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	slog.Info("Sync request received", "method", "POST", "path", "/sync", "remote_addr", r.RemoteAddr)

	if !TokenOk(w, r, model.WTWrite, "") {
		slog.Warn("Sync denied: insufficient write permission", "remote_addr", r.RemoteAddr)
		return
	}
	slog.Info("Synchronizing metadata from disk", "data_dir", config.DataDir)
	store.SyncMetaDataFromDisk()
	slog.Info("Metadata synchronized successfully")
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

	slog.Info("Tag list request received", "method", "GET", "path", "/tag/list", "token_present", clientToken != "", "accessible_tags_count", len(tags), "remote_addr", r.RemoteAddr)

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

	slog.Info("Add tag request received", "method", "POST", "path", "/tag/add", "tag_name", tagName, "remote_addr", r.RemoteAddr)

	if !TokenOk(w, r, model.WTWrite, "") {
		slog.Warn("Add tag denied: insufficient write permission", "tag_name", tagName, "remote_addr", r.RemoteAddr)
		return
	}

	config.AddTagTokenList(tagName)

	slog.Info("Tag added successfully", "tag_name", tagName)

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

	slog.Info("Tag update request received", "method", "POST", "path", "/tag/update", "package_name", pkgName, "new_tag", newTag, "remote_addr", r.RemoteAddr)

	var pkg *model.Package
	pkg, err = store.GetPackage(pkgName)

	if err != nil {
		slog.Info("Tag update failed: package not found", "package_name", pkgName, "remote_addr", r.RemoteAddr)
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("package "+pkgName+" not found"),
		)
		return
	}
	oldTag := pkg.Tag

	slog.Debug("Tag update permission check", "old_tag", oldTag, "new_tag", newTag)

	if !TokenOk(w, r, model.WTWrite, oldTag) || !TokenOk(w, r, model.WTWrite, newTag) {
		slog.Warn("Tag update denied: insufficient write permission", "package_name", pkgName, "old_tag", oldTag, "new_tag", newTag, "remote_addr", r.RemoteAddr)
		return
	}

	err = store.UpdateTag(pkgName, newTag)
	if err != nil {
		slog.Error("Tag update failed: store error", "error", err.Error(), "package_name", pkgName, "old_tag", oldTag, "new_tag", newTag)
		httphelper.SendJSONResponse(
			w,
			http.StatusInternalServerError,
			model.InternalErrorResponse("fail to update the tag:"+err.Error()),
		)
		return
	}

	slog.Info("Tag updated successfully", "package_name", pkgName, "old_tag", oldTag, "new_tag", newTag)
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

	slog.Info("Tag remove request received", "method", "DELETE", "path", "/tag/rm", "tag_name", tag, "remote_addr", r.RemoteAddr)

	if !TokenOk(w, r, model.WTWrite, "") {
		slog.Warn("Tag remove denied: insufficient write permission", "tag_name", tag, "remote_addr", r.RemoteAddr)
		return
	}
	tagList := config.GetTagList()

	if !slices.Contains(tagList, tag) {
		slog.Info("Tag remove failed: tag not found", "tag_name", tag, "available_tags", tagList, "remote_addr", r.RemoteAddr)
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("tag "+tag+" not found"),
		)
		return

	}

	pkgList := store.ListPackagesByTag(tag)
	slog.Debug("Tag removal: packages to be moved", "tag_name", tag, "affected_packages_count", len(pkgList))

	err = store.RemoveTag(tag)
	if err != nil {
		slog.Error("Tag remove failed: store error", "error", err.Error(), "tag_name", tag)
		httphelper.SendJSONResponse(
			w,
			http.StatusInternalServerError,
			model.InternalErrorResponse("error when remove the tag: "+err.Error()),
		)
		return
	}

	err = store.SyncMetaDataToFile()
	if err != nil {
		slog.Error("Tag remove failed: sync error", "error", err.Error(), "tag_name", tag)
		httphelper.SendJSONResponse(
			w,
			http.StatusInternalServerError,
			model.InternalErrorResponse("error when sync the data to disk file"+err.Error()),
		)
		return
	}

	config.DeleteTagTokenList(tag)

	movedPackages := ""
	for i := range pkgList {
		movedPackages += pkgList[i] + "\n"
	}

	slog.Info("Tag removed successfully", "tag_name", tag, "packages_moved_to_temp", len(pkgList))
	successMsg := fmt.Sprintf(
		"remove the tag %s success and packages below has been retaged as temp:\n%v",
		tag, movedPackages,
	)
	httphelper.SendJSONResponse(
		w,
		http.StatusOK,
		model.SuccessfulResponse(successMsg, ""),
	)
}

// reload the server config file from disk which require the admin write token
func reloadHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	slog.Info("Config reload request received", "method", "POST", "path", "/reload", "remote_addr", r.RemoteAddr)

	if !TokenOk(w, r, model.WTWrite, "") {
		slog.Warn("Config reload denied: insufficient write permission", "remote_addr", r.RemoteAddr)
		return
	}

	slog.Info("Reloading server config from disk")

	err := config.InitServerConfig()
	if err != nil {
		slog.Error("Config reload failed", "error", err.Error())
		httphelper.SendJSONResponse(
			w,
			http.StatusInternalServerError,
			model.InternalErrorResponse("error when reload the server config -> "+err.Error()),
		)
		return
	}

	serverConfigPath, _ := config.GetServerConfigPath()
	data := fmt.Sprintf(
		"reload server config file from %s success",
		serverConfigPath,
	)
	slog.Info("Server config reloaded successfully", "config_path", serverConfigPath)
	httphelper.SendJSONResponse(
		w,
		http.StatusOK,
		model.SuccessfulResponse(data, ""),
	)
}

func publicHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := r.URL.Query()
	pkgName := data.Get("name")

	slog.Info("Public request received", "method", "POST", "path", "/public", "package_name", pkgName, "remote_addr", r.RemoteAddr)

	pkg, err := store.GetPackage(pkgName)
	if err != nil {
		slog.Info("Public failed: package not found", "package_name", pkgName, "remote_addr", r.RemoteAddr)
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("pkg "+pkgName+" not found"),
		)
		return
	}

	tag := pkg.Tag

	if !TokenOk(w, r, model.WTWrite, tag) {
		slog.Warn("Public denied: insufficient write permission", "package_name", pkgName, "tag", tag, "remote_addr", r.RemoteAddr)
		return
	}

	slog.Info("Exposing package via Tailscale Funnel", "package_name", pkgName, "tag", tag)
	link, err := exposeSinglePackage(pkgName)
	if err != nil {
		errMsg := fmt.Sprintf(
			"error when public pkg -> %s, error -> %s", pkgName, err.Error(),
		)
		slog.Error("Public failed: funnel error", "error", err.Error(), "package_name", pkgName)
		httphelper.SendJSONResponse(
			w,
			http.StatusInternalServerError,
			model.InternalErrorResponse(errMsg),
		)
		return
	}

	slog.Info("Package made public successfully", "package_name", pkgName, "public_link", link)
	httphelper.SendJSONResponse(
		w,
		http.StatusOK,
		model.SuccessfulResponse(link, ""),
	)
}

func linksHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	slog.Info("Links request received", "method", "GET", "path", "/links", "remote_addr", r.RemoteAddr)

	if !TokenOk(w, r, model.WTRead, "") {
		slog.Warn("Links denied: insufficient read permission", "remote_addr", r.RemoteAddr)
		return
	}

	if len(pkgLinkPool) == 0 {
		slog.Info("No public packages available", "remote_addr", r.RemoteAddr)
		httphelper.SendJSONResponse(
			w,
			http.StatusNotFound,
			model.NotFoundResponse("there is no public package"),
		)
		return
	}

	slog.Info("Returning public links", "links_count", len(pkgLinkPool))
	httphelper.SendJSONResponse(
		w,
		http.StatusOK,
		model.SuccessfulResponse(pkgLinkPool, ""),
	)
}

func privateHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data := r.URL.Query()
	pkgName := data.Get("name")

	slog.Info("Private request received", "method", "POST", "path", "/private", "package_name", pkgName, "remote_addr", r.RemoteAddr)

	if !TokenOk(w, r, model.WTWrite, "") {
		slog.Warn("Private denied: insufficient write permission", "package_name", pkgName, "remote_addr", r.RemoteAddr)
		return
	}

	slog.Info("Making package private", "package_name", pkgName)
	err := privateSinglePackage(pkgName)
	if err != nil {
		slog.Error("Private failed: cannot remove symlink", "error", err.Error(), "package_name", pkgName)
		httphelper.SendJSONResponse(
			w,
			http.StatusInternalServerError,
			model.InternalErrorResponse(err.Error()),
		)
		return
	}

	successMsg := fmt.Sprintf(
		"private the pkg -> %s", pkgName,
	)
	slog.Info("Package made private successfully", "package_name", pkgName)
	httphelper.SendJSONResponse(
		w,
		http.StatusOK,
		model.SuccessfulResponse(successMsg, ""),
	)
	return
}
