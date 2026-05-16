package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"gitee.com/cai-zixiang_hainan/wt/internal/config"
	"gitee.com/cai-zixiang_hainan/wt/internal/model"
	"gitee.com/cai-zixiang_hainan/wt/internal/presets/commonpresets"
)

func searchRequest(pattern string) {
	if pattern == "" {
		fmt.Println("error : package name is required")
		return
	}
	val := url.Values{}
	val.Set("name", pattern)
	apiRsp, err := doRequest(http.MethodGet, "/search", val, config.WTRead, nil)
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
}

func infoRequest(pkgName string) {
	if pkgName == "" {
		fmt.Println("error : package name is required")
		return
	}
	val := url.Values{}
	val.Set("name", pkgName)
	apiRsp, err := doRequest(http.MethodGet, "/search", val, config.WTRead, nil)
	if err != nil {
		fmt.Println("error when search :", err)
		return
	}

	slog.Debug(apiRsp.Message, "func", "searchRequest")

	jsonData, _ := json.Marshal(apiRsp.Data)

	fmt.Println("client jsonData", string(jsonData))

	var results []*model.Package
	err = json.Unmarshal(jsonData, &results)
	if err != nil {
		fmt.Println("fail to unmarshal jsondata to results")
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
}

func doRequest(
	method string,
	endpoint string,
	queryParams url.Values,
	wtMethod config.WTMethod,
	body io.Reader,
) (apiRsp *model.APIResponse, err error) {
	// concate full url
	fullURL := "http://" + config.GetServerAddr(config.WTClient) + endpoint
	if queryParams != nil {
		fullURL += "?" + queryParams.Encode()
	}

	// get token and timeout
	token := config.GetToken(wtMethod)
	timeout := config.GetTimeout(config.WTClient, wtMethod)

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
		headName := config.GetTokenHeadName(config.WTRead)
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
		searchRequest(args[2])
	case "info":
		infoRequest(args[2])
		// case "install":
		// installRequest(args[2])
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
	}
}
