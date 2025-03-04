// GraphQL response processing

package graphql

import (
	"encoding/json"
	"fmt"
	"io"
)

// ***************************************************** //
//                          TYPES                        //
// ***************************************************** //

type gqlResponse struct {
	Data   interface{}
	Errors []gqlError
}

type gqlError struct {
	Message string
	// NOTE: there should be a path (see the spec / the sigma backend error responses)
}

func (e gqlError) Error() string {
	return "graphql: " + e.Message
}

// ***************************************************** //
//                         METHODS                       //
// ***************************************************** //

// unmarshal reads a GraphQL response body into a provided interface.
func unmarshal(body io.Reader, response interface{}) (*gqlResponse, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("Error reading graphql response body: %v", err)
	}

	data := &gqlResponse{
		Data: response,
	}

	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, fmt.Errorf("Error decoding graphql response body: %v", err)
	}

	return data, nil
}
