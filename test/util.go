package test

import (
	"net/http"
	"net/http/httptest"
)

func ServerMock(url string, handler func(http.ResponseWriter, *http.Request)) *httptest.Server {
	h := http.NewServeMux()
	h.HandleFunc(url, handler)
	return httptest.NewServer(h)
}