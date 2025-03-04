// GraphQL requests

package graphql

import (
	"encoding/json"
	"io"
	"net/http"
)

// Request is a GraphQL request
type Request struct {
	query     string
	variables map[string]interface{}
	files     map[string][]file
	header    http.Header
}

// file represents a file to be uploaded
type file struct {
	filename string
	reader   io.Reader
}

// NewRequest creates a new GraphQL request with the given query
func NewRequest(query string) *Request {
	return &Request{
		query:     query,
		variables: make(map[string]interface{}),
		files:     make(map[string][]file),
		header:    make(http.Header),
	}
}

// ***************************************************** //
//                     HELPER METHODS                    //
// ***************************************************** //

// Var sets a variable for the request.
func (r *Request) Var(fieldname string, value interface{}) {
	r.variables[fieldname] = value
}

// Header sets a header for the request.
func (r *Request) Header(key, value string) {
	r.header.Set(key, value)
}

// File sets a file upload for the request.
// Requests containing files will be sent as multipart/form-data.
func (r *Request) File(fieldname, filename string, reader io.Reader) {
	_, exists := r.files[fieldname]

	if exists {
		r.files[fieldname] = append(r.files[fieldname], file{filename, reader})
	} else {
		r.files[fieldname] = []file{file{filename, reader}}
	}
}

// toJSON converts the request to a JSON body containing the query and variables.
func (r *Request) toJSON() ([]byte, error) {
	return json.Marshal(struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}{
		Query:     r.query,
		Variables: r.variables,
	})
}
