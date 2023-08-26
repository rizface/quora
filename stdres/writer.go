package stdres

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

var logger = slog.New(slog.NewJSONHandler(
	os.Stderr,
	&slog.HandlerOptions{
		Level:     slog.LevelWarn,
		AddSource: true,
	},
))

type Response struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	RequestId string      `json:"requestId"`
	Info      string      `json:"info"`
}

func Writer(w http.ResponseWriter, resp Response) {
	w.WriteHeader(resp.Code)

	requestId := w.Header().Get("X-Request-Id")
	if requestId != "" {
		resp.RequestId = requestId
	} else {
		resp.RequestId = uuid.NewString()
		w.Header().Set("X-Request-Id", resp.RequestId)
	}

	if resp.Code >= 500 {
		respB, _ := json.Marshal(resp) //nolint:errcheck

		logger.Error(string(respB))

		resp.Info = "internal server error"
	}

	json.NewEncoder(w).Encode(resp) //nolint:errcheck
}
