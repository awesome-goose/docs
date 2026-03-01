# HTTP Testing

Test your API endpoints.

## Overview

HTTP tests verify that your API endpoints return correct responses.

## Test Setup

### Create Test Server

```go
// test_helper.go
package tests

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "myapp/app"
    "github.com/awesome-goose/goose"
)

func setupTestServer(t *testing.T) *httptest.Server {
    // Create test application
    module := &app.AppModule{}

    // Create HTTP handler
    handler := goose.CreateHandler(module)

    // Create test server
    server := httptest.NewServer(handler)

    t.Cleanup(func() {
        server.Close()
    })

    return server
}
```

### Helper Functions

```go
func makeRequest(t *testing.T, method, url string, body interface{}) *http.Response {
    t.Helper()

    var reqBody io.Reader
    if body != nil {
        jsonBody, _ := json.Marshal(body)
        reqBody = bytes.NewBuffer(jsonBody)
    }

    req, err := http.NewRequest(method, url, reqBody)
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        t.Fatalf("Request failed: %v", err)
    }

    return resp
}

func parseResponse(t *testing.T, resp *http.Response, v interface{}) {
    t.Helper()
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        t.Fatalf("Failed to read response: %v", err)
    }

    if err := json.Unmarshal(body, v); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }
}
```

## Testing Endpoints

### GET Request

```go
func TestGetUsers(t *testing.T) {
    server := setupTestServer(t)

    resp := makeRequest(t, "GET", server.URL+"/api/users", nil)
    defer resp.Body.Close()

    // Assert status code
    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected status 200, got %d", resp.StatusCode)
    }

    // Parse response
    var users []User
    parseResponse(t, resp, &users)

    // Assert response body
    if len(users) == 0 {
        t.Error("expected at least one user")
    }
}
```

### POST Request

```go
func TestCreateUser(t *testing.T) {
    server := setupTestServer(t)

    body := map[string]string{
        "email": "test@example.com",
        "name":  "Test User",
    }

    resp := makeRequest(t, "POST", server.URL+"/api/users", body)
    defer resp.Body.Close()

    // Assert status code
    if resp.StatusCode != http.StatusCreated {
        t.Errorf("expected status 201, got %d", resp.StatusCode)
    }

    // Parse response
    var user User
    parseResponse(t, resp, &user)

    // Assert response
    if user.Email != "test@example.com" {
        t.Errorf("expected email test@example.com, got %s", user.Email)
    }

    if user.ID == "" {
        t.Error("expected user ID to be set")
    }
}
```

### PUT Request

```go
func TestUpdateUser(t *testing.T) {
    server := setupTestServer(t)

    // First create a user
    createBody := map[string]string{
        "email": "original@example.com",
        "name":  "Original Name",
    }
    createResp := makeRequest(t, "POST", server.URL+"/api/users", createBody)

    var createdUser User
    parseResponse(t, createResp, &createdUser)

    // Update the user
    updateBody := map[string]string{
        "name": "Updated Name",
    }

    resp := makeRequest(t, "PUT", server.URL+"/api/users/"+createdUser.ID, updateBody)
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected status 200, got %d", resp.StatusCode)
    }

    var updatedUser User
    parseResponse(t, resp, &updatedUser)

    if updatedUser.Name != "Updated Name" {
        t.Errorf("expected name Updated Name, got %s", updatedUser.Name)
    }
}
```

### DELETE Request

```go
func TestDeleteUser(t *testing.T) {
    server := setupTestServer(t)

    // Create user first
    createBody := map[string]string{
        "email": "delete@example.com",
        "name":  "To Delete",
    }
    createResp := makeRequest(t, "POST", server.URL+"/api/users", createBody)

    var user User
    parseResponse(t, createResp, &user)

    // Delete the user
    resp := makeRequest(t, "DELETE", server.URL+"/api/users/"+user.ID, nil)
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNoContent {
        t.Errorf("expected status 204, got %d", resp.StatusCode)
    }

    // Verify deletion
    getResp := makeRequest(t, "GET", server.URL+"/api/users/"+user.ID, nil)
    if getResp.StatusCode != http.StatusNotFound {
        t.Errorf("expected status 404 after deletion, got %d", getResp.StatusCode)
    }
}
```

## Testing Authentication

### With JWT

```go
func TestProtectedEndpoint(t *testing.T) {
    server := setupTestServer(t)

    t.Run("without token", func(t *testing.T) {
        req, _ := http.NewRequest("GET", server.URL+"/api/protected", nil)

        client := &http.Client{}
        resp, _ := client.Do(req)
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusUnauthorized {
            t.Errorf("expected status 401, got %d", resp.StatusCode)
        }
    })

    t.Run("with valid token", func(t *testing.T) {
        token := generateTestToken("user-123")

        req, _ := http.NewRequest("GET", server.URL+"/api/protected", nil)
        req.Header.Set("Authorization", "Bearer "+token)

        client := &http.Client{}
        resp, _ := client.Do(req)
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            t.Errorf("expected status 200, got %d", resp.StatusCode)
        }
    })

    t.Run("with invalid token", func(t *testing.T) {
        req, _ := http.NewRequest("GET", server.URL+"/api/protected", nil)
        req.Header.Set("Authorization", "Bearer invalid-token")

        client := &http.Client{}
        resp, _ := client.Do(req)
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusUnauthorized {
            t.Errorf("expected status 401, got %d", resp.StatusCode)
        }
    })
}
```

## Testing Validation

```go
func TestCreateUser_Validation(t *testing.T) {
    server := setupTestServer(t)

    tests := []struct {
        name           string
        body           map[string]string
        expectedStatus int
    }{
        {
            name:           "missing email",
            body:           map[string]string{"name": "Test"},
            expectedStatus: http.StatusBadRequest,
        },
        {
            name:           "invalid email",
            body:           map[string]string{"email": "invalid", "name": "Test"},
            expectedStatus: http.StatusBadRequest,
        },
        {
            name:           "missing name",
            body:           map[string]string{"email": "test@example.com"},
            expectedStatus: http.StatusBadRequest,
        },
        {
            name:           "valid input",
            body:           map[string]string{"email": "test@example.com", "name": "Test"},
            expectedStatus: http.StatusCreated,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            resp := makeRequest(t, "POST", server.URL+"/api/users", tc.body)
            defer resp.Body.Close()

            if resp.StatusCode != tc.expectedStatus {
                t.Errorf("expected status %d, got %d", tc.expectedStatus, resp.StatusCode)
            }
        })
    }
}
```

## Testing Error Responses

```go
func TestNotFoundError(t *testing.T) {
    server := setupTestServer(t)

    resp := makeRequest(t, "GET", server.URL+"/api/users/non-existent-id", nil)
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNotFound {
        t.Errorf("expected status 404, got %d", resp.StatusCode)
    }

    var errResp map[string]string
    parseResponse(t, resp, &errResp)

    if errResp["error"] == "" {
        t.Error("expected error message in response")
    }
}
```

## Testing Query Parameters

```go
func TestGetUsersWithPagination(t *testing.T) {
    server := setupTestServer(t)

    // Create some users first
    for i := 0; i < 20; i++ {
        makeRequest(t, "POST", server.URL+"/api/users", map[string]string{
            "email": fmt.Sprintf("user%d@example.com", i),
            "name":  fmt.Sprintf("User %d", i),
        })
    }

    // Test pagination
    resp := makeRequest(t, "GET", server.URL+"/api/users?page=1&limit=10", nil)
    defer resp.Body.Close()

    var users []User
    parseResponse(t, resp, &users)

    if len(users) != 10 {
        t.Errorf("expected 10 users, got %d", len(users))
    }

    // Check pagination headers
    totalCount := resp.Header.Get("X-Total-Count")
    if totalCount == "" {
        t.Error("expected X-Total-Count header")
    }
}
```

## Testing File Uploads

```go
func TestFileUpload(t *testing.T) {
    server := setupTestServer(t)

    // Create multipart form
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)

    // Add file
    part, _ := writer.CreateFormFile("file", "test.txt")
    part.Write([]byte("test file content"))
    writer.Close()

    req, _ := http.NewRequest("POST", server.URL+"/api/upload", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, _ := client.Do(req)
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected status 200, got %d", resp.StatusCode)
    }
}
```

## Testing Response Headers

```go
func TestResponseHeaders(t *testing.T) {
    server := setupTestServer(t)

    resp := makeRequest(t, "GET", server.URL+"/api/users", nil)
    defer resp.Body.Close()

    // Check Content-Type
    contentType := resp.Header.Get("Content-Type")
    if !strings.Contains(contentType, "application/json") {
        t.Errorf("expected JSON content type, got %s", contentType)
    }

    // Check CORS headers
    cors := resp.Header.Get("Access-Control-Allow-Origin")
    if cors == "" {
        t.Error("expected CORS header")
    }
}
```

## Best Practices

1. **Clean up test data** - Reset database between tests
2. **Use table-driven tests** - For testing multiple scenarios
3. **Test all HTTP methods** - GET, POST, PUT, DELETE
4. **Test error responses** - Not just happy paths
5. **Test authentication** - Both with and without tokens
6. **Test validation** - Invalid input scenarios
7. **Check response headers** - Content-Type, CORS, etc.

## Next Steps

- [Integration Testing](integration.md) - Full module tests
- [Unit Testing](unit.md) - Component tests
- [Mocking](mocking.md) - Test doubles
