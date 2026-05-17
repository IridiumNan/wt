package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"gitee.com/cai-zixiang_hainan/wt/internal/config"
	"gitee.com/cai-zixiang_hainan/wt/internal/model"
	"gitee.com/cai-zixiang_hainan/wt/internal/presets/commonpresets"
)

const (
	URLPrefix = "http://"
)

func isEmpty(target string) bool {
	return target == ""
}

func searchRequest(pattern string) (err error) {
	if isEmpty(pattern) {
		err = errors.New("empty pkg name for wt search")
		return
	}

	val := url.Values{}
	val.Set("name", pattern)
	apiRsp, err := doRequest(http.MethodGet, "/search", val, model.WTRead, nil)
	if err != nil {
		fmt.Println("error when search :", err)
		return
	}

	slog.Debug(apiRsp.Message, "func", "searchRequest")

	jsonData, _ := json.Marshal(apiRsp.Data)

	var results []*model.Package
	err = json.Unmarshal(jsonData, &results)
	if err != nil {
		fmt.Println("fail to unmarshal jsondata to results")
		return
	}

	for i := range results {
		fmt.Println(results[i].Name)
	}

	return
}

func infoRequest(pkgName string) (err error) {
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

	fmt.Println("client jsonData", string(jsonData))

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

func installRequest(pkgName string) (err error) {
	if isEmpty(pkgName) {
		err = errors.New("empty pkg name for wt install")
		return
	}

	val := url.Values{}
	val.Set("name", pkgName)

	endpoint := "/install"

	fullURL := URLPrefix + config.GetServerAddr(model.WTClient) + endpoint
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

	slog.Info("Package download successfully",
		"path", localPath,
		"size", written)

	return nil
}

func doRequest(
	method string,
	endpoint string,
	queryParams url.Values,
	wtMethod model.WTMethod,
	body io.Reader,
) (apiRsp *model.APIResponse, err error) {
	// concate full url
	fullURL := URLPrefix + config.GetServerAddr(model.WTClient) + endpoint
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
		headName := config.GetTokenHeadName(model.WTRead)
		req.Header.Set(headName, token)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := httpClient.Do(req)

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

func ClientMain(args []string) {
	if len(args) < 3 {
		fmt.Println(commonpresets.HelpManual)

		fmt.Println("at least 2 parameters")
	}
	command := args[1]

	if command == "--help" || command == "-h" || command == "help" {
		fmt.Println(commonpresets.HelpManual)
		return
	}

	switch command {
	case "--help", "-h", "help":
		fmt.Println(commonpresets.HelpManual)
	case "search":
		err := searchRequest(args[2])
		if err != nil {
			fmt.Println("exec command search fail: ", err)
		}
	case "info":
		err := infoRequest(args[2])
		if err != nil {
			fmt.Println("exec command info fail: ", err)
		}
	case "install":
		err := installRequest(args[2])
		if err != nil {
			fmt.Println("exec command install fail: ", err)
		}
		// case "upload":
		// uploadRequest(args[2])
		// case "replace":
		// 	if len(args) < 4 {
		// 		fmt.Println("Usage: wt replace <package name> <path to your new package>")
		// 		return
		// 	}
		// 	replaceRequest(args[2], args[3])
		// case "mv":
		// 	if len(args) < 4 {
		// 		fmt.Println("Usage: wt mv <package name> <new package name>")
		// 		return
		// 	}
		// 	mvRequest(args[2], args[3])
		// case "rm":
		// 	rmRequest(args[2])
	default:
		fmt.Println(commonpresets.HelpManual)
	}
}
