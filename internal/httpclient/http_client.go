package httpclient

import (
	"fmt"
	"github.com/badrchoubai/atoolforparsingandrunninghttpfiles/internal/logging"
	"github.com/badrchoubai/atoolforparsingandrunninghttpfiles/internal/parser"
	"net/http"
	"time"
)

var _ HTTPClient = (*httpClient)(nil)

type HTTPClient interface {
	DoRequest(request *parser.HTTPRequest) (response *http.Response, err error)
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

func (h *httpClient) DoRequest(request *parser.HTTPRequest) (*http.Response, error) {
	switch request.Method {
	case http.MethodGet:
		return h.get(request)
	default:
		return nil, fmt.Errorf("unsupported method %s", request.Method)
	}
}

func (h *httpClient) get(request *parser.HTTPRequest) (*http.Response, error) {
	resp, err := h.client.Get(request.URL)
	if err != nil {
		h.logger.Error(fmt.Sprintf("GET request failed: %v", err))
		return nil, err
	}

	h.logger.Info(fmt.Sprintf("GET request succeeded with status: %s", resp.Status))
	return resp, nil

}
