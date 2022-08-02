package myhttp

import "net/http"

const (
	CONTENT_TYPE_JS = "application/javascript; charset=utf-8"
)

//HTTP响应文本内容
func ResponseString(w http.ResponseWriter, body string, content_type string) {
	w.Header().Set("content-type", content_type)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}
