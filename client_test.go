package graphql

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Needed because multipart dictates that line endings should be \r\n
func normalizeLineEndings(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// TestJson tests the production of JSON
func TestJson(t *testing.T) {
	request := NewRequest(`query User($id: ID!) { user(id: $id) { name } }`)
	request.Var("id", 1)

	client := NewClient("")

	result, err := client.requestJson(request)
	if err != nil {
		t.Fatal(err)
	}

	// Read the request body
	contentBytes, err := io.ReadAll(result.Body)
	expected := `{"query":"query User($id: ID!) { user(id: $id) { name } }","variables":{"id":1}}`
	assert.Equal(t, expected, string(contentBytes))
}

// TestMultipart tests the production of multipart requests.
func TestMultipart(t *testing.T) {
	// https://github.com/jaydenseric/graphql-multipart-request-spec#operations
	request := NewRequest(`mutation ($file: Upload!) { singleUpload(file: $file) { id } }`)

	request.File("file", "a.txt", strings.NewReader("Alpha file content.\n"))

	client := NewClient("")

	result, err := client.requestMultipart(request)
	if err != nil {
		t.Fatal(err)
	}

	// Read the request body
	contentBytes, err := io.ReadAll(result.Body)

	// Parse the boundary string
	content := normalizeLineEndings(string(contentBytes))
	boundary := strings.Split(content, "\n")[0][2:]

	expectedTemplate := `--%s
Content-Disposition: form-data; name="operations"

{"query":"mutation ($file: Upload!) { singleUpload(file: $file) { id } }","variables":{"file":null}}
--%s
Content-Disposition: form-data; name="map"

{"0":["variables.file"]}
--%s
Content-Disposition: form-data; name="0"; filename="a.txt"
Content-Type: text/plain; charset=utf-8

Alpha file content.

--%s--
`
	expected := fmt.Sprintf(expectedTemplate, boundary, boundary, boundary, boundary)

	assert.Equal(t, expected, content)
}

// TestUnmarshal tests the unmarshalling of different types in a GraphQL request.
func TestUnmarshal(t *testing.T) {
	response := `
  {
    "data": {
      "string": "string",
      "int": 1,
      "float": 1.1,
      "bool": true
    }
  }
  `

	data := struct {
		String string
		Int    int
		Float  float64
		Bool   bool
	}{}

	_, err := unmarshal(strings.NewReader(response), &data)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, data.String, "string")
	assert.Equal(t, data.Int, 1)
	assert.Equal(t, data.Float, 1.1)
	assert.Equal(t, data.Bool, true)
}

// TestUnmarshalErrors tests the unmarshalling of errors in a GraphQL request.
func TestUnmarshalErrors(t *testing.T) {
	response := `
  {
    "errors": [
      {
        "message": "User not found."
      },
      {
        "message": "Permission denied."
      }
    ]
  }
  `

	data := struct {
	}{}

	result, err := unmarshal(strings.NewReader(response), &data)

	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, result.Errors, 2)
	assert.Equal(t, "User not found.", result.Errors[0].Message)
	assert.Equal(t, "Permission denied.", result.Errors[1].Message)
}
