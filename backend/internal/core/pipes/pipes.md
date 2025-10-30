# Pipes

The pipes for this project only serve the validation before the request reaches the handler
## Overview

In Gin, the `c.Next()` function is used to continue processing the request down the middleware chain. You can use this function in your interceptor to:

- **Before `c.Next()`:** Perform actions before the request is processed by the subsequent handlers or middleware.
- **After `c.Next()`:** Modify the response or perform additional actions after the request has been processed.

## Example Usage

Here's an example showing how to use a pipe to validate the request body before it reaches the handler

```go
func BodyValidator[T any](c *gin.Context) {
	var dto T // Data Transfer Object

	// Bind JSON to the DTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		// Append the error to the context without aborting the request
		c.Error(fmt.Errorf("Invalid Input"))
		// Set the status code to 400
		c.Set("statusCode", http.StatusBadRequest)
	}

	// Validate the DTO using the validator instance
	if err := validate.Struct(dto); err != nil {
		// Append the validation error to the context without aborting the request
		c.Error(fmt.Errorf("Invalid Input"))

		// Set the status code to 400
		c.Set("statusCode", http.StatusBadRequest)
	}

	// Proceed to the next middleware or handler
	c.Next()
}
