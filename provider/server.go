package provider

import (
	"net/http"
	"os"
)

func ProviderServer(h http.Handler) *http.Server {
	return &http.Server{
		Addr:    os.Getenv("APP_PORT"),
		Handler: h,
	}
}
