// JSON content type requests.

package graphql

import (
	"bytes"
	"fmt"
	"net/http"
)

// requestJson creates a HTTP request from a GraphQL request using JSON content type.
func (c *Client) requestJson(request *Request) (*http.Request, error) {
	requestBody, err := request.toJSON()
	if err != nil {
		return nil, fmt.Errorf("Error marshaling the graphql request to json: %v", err)
	}

	r, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("Error creating graphql request: %v", err)
	}

	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	r.Header.Set("Accept", "application/json; charset=utf-8")

	for key, values := range request.header {
		for _, value := range values {
			r.Header.Add(key, value)
		}
	}

	return r, nil
}
