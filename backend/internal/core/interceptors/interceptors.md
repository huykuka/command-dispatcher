# Gin Interceptor Example

The middlewares for this project serve the actions should be executed  before and after  the request reaches the handler

## Overview

In Gin, the `c.Next()` function is used to continue processing the request down the middleware chain. You can use this function in your interceptor to:

- **Before `c.Next()`:** Perform actions before the request is processed by the subsequent handlers or middleware.
- **After `c.Next()`:** Modify the response or perform additional actions after the request has been processed.

## Example Usage

Here's an example showing how to use an interceptor to log the request details before processing and modify the response afterward:

```go
func ExampleInterceptor() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Do something before the request is handled
        fmt.Println("Before request handling")

        // Continue processing the request
        c.Next()

        // Do something after the request is handled
        fmt.Println("After request handling")

        // Example: Modifying the response
        if c.Writer.Status() == http.StatusOK {
            // Modify the response if the status is OK
            c.JSON(http.StatusOK, gin.H{
                "message": "This is a modified response",
            })
        }
    }
}
