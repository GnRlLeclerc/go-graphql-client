// Multipat content type requests

package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
)

func (c *Client) requestMultipart(request *Request) (*http.Request, error) {
	var requestBody bytes.Buffer
	fileMap := make(map[string][]string)
	writer := multipart.NewWriter(&requestBody)

	// Add file null variables to the request for each file to be uploaded
	filecount := 0
	for fieldname, files := range request.files {
		if len(files) > 1 {
			// Slice of nil
			request.Var(fieldname, make([]interface{}, len(files)))
			for i := range files {
				fileMap[strconv.Itoa(filecount)] = []string{fmt.Sprintf("variables.%s.%d", fieldname, i)}
				filecount += 1
			}
		} else {
			request.Var(fieldname, nil)
			fileMap[strconv.Itoa(filecount)] = []string{fmt.Sprintf("variables.%s", fieldname)}
			filecount += 1
		}
	}

	query, err := request.toJSON()
	if err != nil {
		return nil, fmt.Errorf("Error marshaling the graphql request to json: %v", err)
	}

	// Write the GraphQL query
	writer.WriteField("operations", string(query))

	// Write the file map
	res, err := json.Marshal(fileMap)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling the file map: %v", err)
	}
	writer.WriteField("map", string(res))

	// TODO: write the files

	// DEBUG
	println(requestBody.String())

	return nil, nil
}
