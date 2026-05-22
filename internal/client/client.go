package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"gitee.com/cai-zixiang_hainan/wt/internal/config"
	"gitee.com/cai-zixiang_hainan/wt/internal/model"
)

func isEmpty(target string) bool {
	return target == ""
}

func SearchRequest(pattern string) (err error) {
	if isEmpty(pattern) {
		err = errors.New("empty pkg name for wt search")
		return
	}

	val := url.Values{}
	val.Set("name", pattern)
	apiRsp, err := doRequest(http.MethodGet, "/search", val, model.WTRead, nil)
	if err != nil {
		return
	}

	slog.Debug(apiRsp.Message, "func", "searchRequest")

	jsonData, _ := json.Marshal(apiRsp.Data)

	var results []*model.Package
	err = json.Unmarshal(jsonData, &results)
	if err != nil {
		return
	}

	for i := range results {
		fmt.Println(results[i].Name)
	}

	return
}

func InfoRequest(pkgName string) (err error) {
	if isEmpty(pkgName) {
		err = errors.New("empty pkg name for wt info")
		return
	}

	val := url.Values{}
	val.Set("name", pkgName)
	apiRsp, err := doRequest(http.MethodGet, "/search", val, model.WTRead, nil)
	if err != nil {
		return
	}

	slog.Debug(apiRsp.Message, "func", "searchRequest")

	jsonData, _ := json.Marshal(apiRsp.Data)

	var results []*model.Package
	err = json.Unmarshal(jsonData, &results)
	if err != nil {
		return
	}

	fmt.Println(results)
	pkgNum := len(results)
	if pkgNum == 1 {
		fmt.Println("------------ pkg information ------------")
		fmt.Println(results[0].Info())
		fmt.Println("-----------------------------------------")
		return
	}

	fmt.Println(pkgNum, " package found")
	for i := range results {
		fmt.Println(">>> package ", i, " <<<")
		fmt.Println("------------ pkg information ------------")
		fmt.Println(results[i].Info())
		fmt.Println("-----------------------------------------")
	}

	return
}

func InstallRequest(pkgName string) (err error) {
	if isEmpty(pkgName) {
		err = errors.New("empty pkg name for wt install")
		return
	}

	val := url.Values{}
	val.Set("name", pkgName)

	endpoint := "/install"

	fullURL := config.GetServerAddr(model.WTClient) + endpoint
	fullURL += "?" + val.Encode()

	token := config.GetToken(model.WTInstall)
	timeout := config.GetTimeout(model.WTClient, model.WTInstall)

	httpClient := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return fmt.Errorf("fail to create request: %w", err)
	}

	headName := config.GetTokenHeadName(model.WTInstall)
	req.Header.Set(headName, token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request fail : %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server return status : %d\n error: %w", resp.StatusCode, body)
	}

	localPath := pkgName

	outFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("fail to create %s here\nerror : %w", localPath, err)
	}
	defer outFile.Close()

	written, err := io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("crash when installing pkg: %s\nsaved size: %d\nerror: %w", localPath, written, err)
	}

	// slog.Info("Package download successfully",
	// 	"path", localPath,
	// 	"size", written)
	fmt.Printf("Install the package: %s, size: %d", pkgName, written)

	return nil
}

func ListRequest(targetTag string) (err error) {
	if isEmpty(targetTag) {
		err = errors.New("empty target tag for wt list")
		return
	}

	val := url.Values{}
	val.Set("tag", targetTag)
	apiRsp, err := doRequest(http.MethodGet, "/list", val, model.WTRead, nil)
	if err != nil {
		return
	}

	slog.Debug(apiRsp.Message, "func", "listRequest")

	slog.Debug("received data ", "data", apiRsp.Data)

	var nameList []string
	jsonData, _ := json.Marshal(apiRsp.Data)

	err = json.Unmarshal(jsonData, &nameList)
	if err != nil {
		return
	}

	for i := range nameList {
		fmt.Println(nameList[i])
	}

	return
}

func UploadRequest(filePath string, pkgName string) (err error) {
	if isEmpty(filePath) {
		return errors.New("can not assign empty file path")
	}
	if isEmpty(pkgName) {
		pkgName = filepath.Base(filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("fail to open local file : %s\nerr : %w", filePath, err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	err = writer.WriteField("name", pkgName)
	if err != nil {
		return fmt.Errorf("fail to write name field: %s\nerr: %w", pkgName, err)
	}

	part, err := writer.CreateFormFile("file", pkgName)
	if err != nil {
		return fmt.Errorf("fail to create form file : %w", err)
	}

	if _, err = io.Copy(part, file); err != nil {
		return fmt.Errorf("fail to copy file content : %w", err)
	}

	writer.Close()

	fullURL := config.GetServerAddr(model.WTClient) + "/upload"
	token := config.GetToken(model.WTWrite)
	timeout := config.GetTimeout(model.WTClient, model.WTWrite)

	req, err := http.NewRequest(http.MethodPost, fullURL, body)
	if err != nil {
		return fmt.Errorf("fail to create request %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if token != "" {
		req.Header.Set(config.GetTokenHeadName(model.WTWrite), token)
	}

	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request fail : %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("upload failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	fmt.Printf("upload file success, localFile -> %s, pkgName -> %s", filePath, pkgName)
	return nil
}

func MvRequest(oldName string, newName string) (err error) {
	if isEmpty(oldName) || isEmpty(newName) {
		return fmt.Errorf("missing old name or new name")
	}

	fmt.Printf("send request for mv %s to %s (server : %s)", oldName, newName, config.GetServerAddr(model.WTClient))
	val := url.Values{}
	val.Set("old_name", oldName)
	val.Set("new_name", newName)
	apiRsp, err := doRequest(http.MethodPost, "/mv", val, model.WTWrite, nil)
	if err != nil {
		return
	}

	slog.Debug(apiRsp.Message, "func", "mvRequest")

	jsonData, _ := json.Marshal(apiRsp.Data)

	slog.Debug("client jsonData", "data", string(jsonData))

	fmt.Println("Status Code ", apiRsp.Code)
	fmt.Println(string(jsonData))

	return
}

func SyncRequest() (err error) {
	fmt.Printf("send request for sync file info from disk (server: %s)", config.GetServerAddr(model.WTServer))
	val := url.Values{}
	apiRsp, err := doRequest(http.MethodPut, "/sync", val, model.WTWrite, nil)
	if err != nil {
		return
	}

	jsonData, _ := json.Marshal(apiRsp.Data)

	fmt.Println("Status Code ", apiRsp.Code)
	fmt.Println(string(jsonData))

	return nil
}

func RmRequest(pkgName string) (err error) {
	if isEmpty(pkgName) {
		return fmt.Errorf("missing pkg name")
	}

	fmt.Printf("send request for rm package : %s (server: %s)", pkgName, config.GetServerAddr(model.WTServer))
	val := url.Values{}
	val.Set("name", pkgName)
	apiRsp, err := doRequest(http.MethodDelete, "/rm", val, model.WTWrite, nil)
	if err != nil {
		return
	}

	slog.Debug(apiRsp.Message, "func", "rmRequest")

	jsonData, _ := json.Marshal(apiRsp.Data)

	fmt.Println("Status Code ", apiRsp.Code)
	fmt.Println(string(jsonData))

	return
}

func TagListRequest() (err error) {
	val := url.Values{}
	apiRsp, err := doRequest(http.MethodGet, "/tag/list", val, model.WTRead, nil)
	if err != nil {
		return
	}

	jsonData, _ := json.Marshal(apiRsp.Data)

	var tags []string
	// return data : tags []string
	err = json.Unmarshal(jsonData, &tags)
	if err != nil {
		return
	}

	for i := range tags {
		fmt.Println(tags[i])
	}

	return
}

func AddTagRequest(tagName string) (err error) {
	val := url.Values{}
	val.Set("tag", tagName)

	fmt.Println("send request for add new tag:", tagName)

	var apiRsp *model.APIResponse
	apiRsp, err = doRequest(http.MethodPost, "/tag/add", val, model.WTWrite, nil)
	if err != nil {
		return
	}

	jsonData, _ := json.Marshal(apiRsp.Data)

	var res string

	err = json.Unmarshal(jsonData, &res)

	fmt.Println(res)

	return
}

func UpdateTagRequest(pkgName string, newTag string) (err error) {
	val := url.Values{}
	val.Set("name", pkgName)
	val.Set("new_tag", newTag)

	fmt.Printf("send request for update the package : %s to new tag : %s (server : %s)", pkgName, newTag, config.GetServerAddr(model.WTClient))

	var apiRsp *model.APIResponse
	apiRsp, err = doRequest(http.MethodPut, "/tag/update", val, model.WTWrite, nil)
	if err != nil {
		return
	}

	if apiRsp.Code != http.StatusOK {
		return errors.New(apiRsp.Error)
	}
	data := apiRsp.Data
	fmt.Println(data)

	return
}

func TagRmRequest(tagName string) (err error) {
	val := url.Values{}
	val.Set("tag", tagName)

	fmt.Println("send request for rm the tag ", tagName)

	var apiRsp *model.APIResponse
	apiRsp, err = doRequest(http.MethodPost, "/tag/rm", val, model.WTWrite, nil)
	if err != nil {
		return
	}

	if apiRsp.Code != http.StatusOK {
		return errors.New(apiRsp.Error)
	}

	data := apiRsp.Data
	fmt.Println(data)

	return
}

func ReloadRequest() (err error) {
	val := url.Values{}

	fmt.Println("send request for reload the server config")
	var apiRsp *model.APIResponse
	apiRsp, err = doRequest(http.MethodPost, "/reload", val, model.WTWrite, nil)
	if err != nil {
		return
	}

	if apiRsp.Code != http.StatusOK {
		return errors.New(apiRsp.Error)
	}

	data := apiRsp.Data

	fmt.Println(data)

	return
}

func PublicRequest(pkgName string) (err error) {
	val := url.Values{}
	val.Set("name", pkgName)

	fmt.Println("send request for public the pkg ->", pkgName)

	var apiRsp *model.APIResponse
	apiRsp, err = doRequest(http.MethodPost, "/public", val, model.WTWrite, nil)
	if err != nil {
		return
	}

	if apiRsp.Code != http.StatusOK {
		return errors.New(apiRsp.Error)
	}
	data := apiRsp.Data
	fmt.Println("public link -> ", data)

	return
}

func printLinks(links []string) {
	fmt.Println("--------------------- available links ------------------------")
	for i := range links {
		fmt.Println(links[i])
	}
	fmt.Println("--------------------------------------------------------------")
}

func LinksRequest() (err error) {
	val := url.Values{}
	fmt.Println("send request for exist links")

	var apiRsp *model.APIResponse
	apiRsp, err = doRequest(http.MethodGet, "/links", val, model.WTRead, nil)
	if err != nil {
		return
	}

	if apiRsp.Code != http.StatusOK {
		return errors.New(apiRsp.Error)
	}

	data := apiRsp.Data
	if links, ok := data.([]string); !ok {
		printLinks(links)
		return
	}

	return errors.New("can't read links")
}

// doRequest : this function will choose the corect token
func doRequest(
	method string,
	endpoint string,
	queryParams url.Values,
	wtMethod model.WTMethod,
	body io.Reader,
) (apiRsp *model.APIResponse, err error) {
	// concate full url
	fullURL := config.GetServerAddr(model.WTClient) + endpoint
	if queryParams != nil {
		fullURL += "?" + queryParams.Encode()
	}

	// get token and timeout
	token := config.GetToken(wtMethod)
	timeout := config.GetTimeout(model.WTClient, wtMethod)

	// construrct the Client
	httpClient := &http.Client{
		Timeout: timeout,
	}

	// create request
	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("fail to create request : %w", err)
	}

	// set the token to Header
	if token != "" {
		headName := config.GetTokenHeadName(wtMethod)
		req.Header.Set(headName, token)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request err :\nreponse : %v\nerr : %w", resp, err)
	}

	// Must close the response.Body
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Println("fail to close response body: ", err)
		}
	}()

	// read respond body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fail to read response: %w", err)
	}

	var apiResp model.APIResponse
	err = json.Unmarshal(respBody, &apiResp)
	slog.Debug(string(respBody))
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal response (status: %d) %s - %w",
			resp.StatusCode, string(respBody), err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("server error : %s", apiResp.Error)
	}

	return &apiResp, nil
}
