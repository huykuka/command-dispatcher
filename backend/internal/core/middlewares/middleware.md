# Gin Middleware 

The middlewares for this project only serve the actions should be executed before the request reaches the handler
## Overview

In Gin, the `c.Next()` function is used to continue processing the request down the middleware chain. You can use this function in your interceptor to:

- **Before `c.Next()`:** Perform actions before the request is processed by the subsequent handlers or middleware.
- **After `c.Next()`:** Modify the response or perform additional actions after the request has been processed.

## Example Usage

Here's an example showing how to use an middleware to log the request details before processing and modify the response afterward:

```go
func ExampleMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Continue processing the request
        c.Next()
    }
}
