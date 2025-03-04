# Go Graphql Client

A simple Graphql client for Golang, inspired by the old [Machinebox Graphql Client](https://github.com/machinebox/graphql).

Support multipart file upload as specified in the [graphql-multipart-request-spec](https://github.com/jaydenseric/graphql-multipart-request-spec).

All requests use Json transport by default, unless files are submitted in which case multipart transport is used.
This client was made for testing purposes, do not use it in production.

## Installation

Run:

```bash
go get github.com/GnRlLeclerc/go-graphql-client
```

## Usage

```go
import "context"

// Create a client (thread-safe)
client := graphql.NewClient("https://localhost:8080/graphql")

// Set client cookies
client.AddCookie("https://localhost:8080/graphql", "auth", "0000")

// Create a request
req := graphql.NewRequest(`
    query ($key: String!) {
        items (id:$key) {
            field1
            field2
            field3
        }
    }
`)

// Set request variables
req.Var("key", "value")

// Set header fields
req.Header("Cache-Control", "no-cache")

// Use any context
ctx := context.Background()

// Run the request and scan the response into a struct
var reponse ResponseStruct
if err := client.Run(ctx, req, &response); err != nil {
    log.Fatal(err)
}

client.ClearCookies()  // Clear cookies if needed (to reuse the client)
```

### Uploading Files

This client supports uploading one or multiple files at the same time.
If including files in a request, multipart transport will be used.

```go
import "context"

// Create a client (thread-safe)
client := graphql.NewClient("https://localhost:8080/graphql")

// Create a request
req := graphql.NewRequest(`
    mutation ($file: Upload!) {
        singleUpload(file: $file) {
            id
        }
    }
`)

// Open a file
file := os.Open("a.txt")
defer file.Close()

req.File("file", "a.txt", file)

// Use any context
ctx := context.Background()

// Run the request and scan the response into a struct
var reponse ResponseStruct
if err := client.Run(ctx, req, &response); err != nil {
    log.Fatal(err)
}
```
