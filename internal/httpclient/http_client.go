package httpclient

import (
	"github.com/badrchoubai/atoolforparsingandrunninghttpfiles/internal/logging"
	"github.com/badrchoubai/atoolforparsingandrunninghttpfiles/internal/parser"
	"net/http"
	"time"
)

var _ HTTPClient = (*httpClient)(nil)

type HTTPClient interface {
	Get(request *parser.HTTPRequest) (resp *http.Response, err error)
}

type httpClient struct {
	client *http.Client
	logger *logging.Logger
}

func NewHTTPClient(logger *logging.Logger) HTTPClient {
	return &httpClient{
		logger: logger,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (h *httpClient) Get(request *parser.HTTPRequest) (resp *http.Response, err error) {
	h.logger.Info(request.Method, "url", request.URL)

	return h.client.Get(request.URL)
}
