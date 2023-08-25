package stdres

import (
	"encoding/json"
	"net/http"
	"os"

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
	const (
		InternalServerErrorJson = `{"code": 500, "info": "internal server error"}`
	)

	var internalServerErrorResponse = func() {
		logger.Error("failed to write response")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(InternalServerErrorJson)) //nolint:errcheck
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Code)

	respB, err := json.Marshal(resp)
	if err != nil {
		internalServerErrorResponse()
	}

	if resp.Code >= 500 {
		logger.Error(string(respB))
	}

	_, err = w.Write(respB)
	if err != nil {
		internalServerErrorResponse()
	}
}
