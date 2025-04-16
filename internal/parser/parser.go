// Package parser contains code for parsing the `.http` file format
package parser

import (
	"bufio"
	"os"
	"strings"
)

var _ Parser = (*HttpFileParser)(nil)

type Parser interface {
	Parse(filepath string) (bool, error)
}

type HTTPRequest struct {
	Description string
	Method      string
	URL         string
}

type HttpFileParser struct {
	Requests []*HTTPRequest
}

// Parse scans a given `.http` file and appends valid results to `HttpFileParser.requests`
// If the list of requests was built Parse returns `true`, otherwise `false` and a resulting error
func (h *HttpFileParser) Parse(filepath string) (bool, error) {
	file, err := os.Open(filepath)
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	if err != nil {
		return false, err
	}

	var request *HTTPRequest
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if description, found := strings.CutPrefix(line, "###"); found {
			description = strings.TrimSpace(description)
			if description != "" {
				if request != nil {
					h.Requests = append(h.Requests, request)
				}
				request = &HTTPRequest{
					Description: description,
				}
				continue
			}
		}

		if request != nil && request.Method == "" && request.URL == "" {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				request.Method = parts[0]
				request.URL = parts[1]
			}
		}
	}

	if request != nil {
		h.Requests = append(h.Requests, request)
	}

	if len(h.Requests) > 0 {
		return true, nil
	} else {
		return false, err
	}
}

func (h *HttpFileParser) addRequest(request *HTTPRequest) {
	h.Requests = append(h.Requests, request)
}
