package auth_hub

import (
	"net/http"
	"time"
)

// Worker 请求机器
type Worker struct {
	client *http.Client
}

// NewWork 禁止重定向的HTTPClient
func NewWork() *Worker {
	return &Worker{client: &http.Client{Timeout: 30 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }}}
}
