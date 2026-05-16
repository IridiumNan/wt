package httphelper

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"gitee.com/cai-zixiang_hainan/wt/internal/model"
)

// SendJSONResponse : receive the APIResponse struct and convert to []byte then send
func SendJSONResponse(w http.ResponseWriter, statusCode int, response *model.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonByte, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	slog.Debug("send data" + string(jsonByte))
	w.Write([]byte(jsonByte))
}
