# HTTP Testing

Test your API endpoints using the `HTTPTest` utility from goose/testing.

## Import

```go
import (
    "net/http"
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)
```

## Overview

HTTP tests verify that your API endpoints return correct responses. The `HTTPTest` helper provides a fluent API for building requests and asserting responses.

## Test Setup

### Using HTTPTest

```go
// test_helper.go
package tests

import (
    "net/http"
    "testing"

    "myapp/app"
    goosetest "github.com/awesome-goose/goose/testing"
)

func setupTestServer(t *testing.T) *goosetest.HTTPTest {
    // Create your HTTP handler
    handler := createHandler(&app.AppModule{})

    // Create HTTPTest with the handler
    ht := goosetest.NewHTTPTest(t, handler)

    // Start a test server (for real HTTP requests)
    ht.WithServer()

    t.Cleanup(func() {
        ht.Close()
    })

    return ht
}
```

## Testing Endpoints

### GET Request

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestGetUsers(t *testing.T) {
    ht := setupTestServer(t)

    // Fluent API for GET request
    ht.GET("/api/users").
        Do().
        ExpectOK()
}

func TestGetUserByID(t *testing.T) {
    ht := setupTestServer(t)

    var user User
    ht.GET("/api/users/123").
        Do().
        ExpectOK().
        JSON(&user)

    ht.T.Expect(user.ID).ToEqual("123")
}
```

### POST Request

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestCreateUser(t *testing.T) {
    ht := setupTestServer(t)

    body := map[string]string{
        "email": "test@example.com",
        "name":  "Test User",
    }

    var user User
    ht.POST("/api/users").
        WithJSON(body).
        Do().
        ExpectCreated().
        JSON(&user)

    ht.T.Expect(user.Email).ToEqual("test@example.com")
    ht.T.Expect(user.ID).Not().ToEqual("")
}
```

### PUT Request

```go
func TestUpdateUser(t *testing.T) {
    ht := setupTestServer(t)

    // First create a user
    var createdUser User
    ht.POST("/api/users").
        WithJSON(map[string]string{"email": "original@example.com", "name": "Original"}).
        Do().
        ExpectCreated().
        JSON(&createdUser)

    // Update the user
    var updatedUser User
    ht.PUT("/api/users/" + createdUser.ID).
        WithJSON(map[string]string{"name": "Updated Name"}).
        Do().
        ExpectOK().
        JSON(&updatedUser)

    ht.T.Expect(updatedUser.Name).ToEqual("Updated Name")
}
```

### DELETE Request

```go
func TestDeleteUser(t *testing.T) {
    ht := setupTestServer(t)

    // Create user first
    var user User
    ht.POST("/api/users").
        WithJSON(map[string]string{"email": "delete@example.com", "name": "To Delete"}).
        Do().
        ExpectCreated().
        JSON(&user)

    // Delete the user
    ht.DELETE("/api/users/" + user.ID).
        Do().
        ExpectNoContent()

    // Verify deletion
    ht.GET("/api/users/" + user.ID).
        Do().
        ExpectNotFound()
}
```

### PATCH Request

```go
func TestPatchUser(t *testing.T) {
    ht := setupTestServer(t)

    var user User
    ht.PATCH("/api/users/123").
        WithJSON(map[string]string{"name": "Patched Name"}).
        Do().
        ExpectOK().
        JSON(&user)

    ht.T.Expect(user.Name).ToEqual("Patched Name")
}
```

## Testing Authentication

### With Bearer Token

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestProtectedEndpoint(t *testing.T) {
    ht := setupTestServer(t)

    t.Run("without token", func(t *testing.T) {
        ht.GET("/api/protected").
            Do().
            ExpectUnauthorized()
    })

    t.Run("with valid token", func(t *testing.T) {
        token := generateTestToken("user-123")

        ht.GET("/api/protected").
            WithBearer(token).
            Do().
            ExpectOK()
    })

    t.Run("with invalid token", func(t *testing.T) {
        ht.GET("/api/protected").
            WithBearer("invalid-token").
            Do().
            ExpectUnauthorized()
    })
}
```

### With Basic Auth

```go
func TestBasicAuth(t *testing.T) {
    ht := setupTestServer(t)

    t.Run("with valid credentials", func(t *testing.T) {
        ht.GET("/api/admin").
            WithBasicAuth("admin", "password").
            Do().
            ExpectOK()
    })

    t.Run("with invalid credentials", func(t *testing.T) {
        ht.GET("/api/admin").
            WithBasicAuth("admin", "wrong").
            Do().
            ExpectUnauthorized()
    })
}
```

## Testing Validation

```go
import (
    "net/http"
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestCreateUser_Validation(t *testing.T) {
    ht := setupTestServer(t)

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
            ht.POST("/api/users").
                WithJSON(tc.body).
                Do().
                ExpectStatus(tc.expectedStatus)
        })
    }
}
```

## Testing Error Responses

```go
func TestNotFoundError(t *testing.T) {
    ht := setupTestServer(t)

    var errResp map[string]string
    ht.GET("/api/users/non-existent-id").
        Do().
        ExpectNotFound().
        JSON(&errResp)

    ht.T.Expect(errResp["error"]).Not().ToEqual("")
}

func TestInternalServerError(t *testing.T) {
    ht := setupTestServer(t)

    ht.GET("/api/trigger-error").
        Do().
        ExpectInternalServerError()
}

func TestForbiddenAccess(t *testing.T) {
    ht := setupTestServer(t)

    ht.DELETE("/api/admin/users/123").
        WithBearer(regularUserToken).
        Do().
        ExpectForbidden()
}
```

## Testing Query Parameters

```go
func TestGetUsersWithPagination(t *testing.T) {
    ht := setupTestServer(t)

    var users []User
    resp := ht.GET("/api/users").
        WithQuery("page", "1").
        WithQuery("limit", "10").
        Do().
        ExpectOK().
        JSON(&users)

    ht.T.Expect(len(users)).ToBeLessThanOrEqual(10)

    // Check pagination headers
    resp.ExpectHeader("X-Total-Count", "20")
}
```

## Testing Headers

```go
func TestCustomHeaders(t *testing.T) {
    ht := setupTestServer(t)

    ht.GET("/api/data").
        WithHeader("X-Custom-Header", "custom-value").
        WithHeaders(map[string]string{
            "X-Request-ID": "12345",
            "Accept":       "application/json",
        }).
        Do().
        ExpectOK().
        ExpectHeader("Content-Type", "application/json")
}
```

## Response Assertions Reference

The `HTTPResponse` provides these assertion methods:

```go
import (
    "testing"

    goosetest "github.com/awesome-goose/goose/testing"
)

func TestResponseAssertions(t *testing.T) {
    ht := setupTestServer(t)

    resp := ht.GET("/api/users").Do()

    // Status assertions
    resp.ExpectOK()                    // 200
    resp.ExpectCreated()               // 201
    resp.ExpectNoContent()             // 204
    resp.ExpectBadRequest()            // 400
    resp.ExpectUnauthorized()          // 401
    resp.ExpectForbidden()             // 403
    resp.ExpectNotFound()              // 404
    resp.ExpectInternalServerError()   // 500
    resp.ExpectStatus(202)             // Custom status

    // Header assertions
    resp.ExpectHeader("Content-Type", "application/json")
    resp.ExpectContentType("application/json")

    // Body access
    body := resp.Body()          // []byte
    bodyStr := resp.BodyString() // string

    // JSON parsing
    var data map[string]any
    resp.JSON(&data)
}
```

## Next Steps

- [Integration Testing](integration.md) - Test module interactions
- [E2E Testing](e2e.md) - Test full user flows
