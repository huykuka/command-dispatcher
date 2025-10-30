# Use Alpine as the base image
FROM alpine:latest

# # Install necessary runtime dependencies
RUN apk add --no-cache postgresql-client

# Set the working directory
WORKDIR /app

# Copy the backend build output from the build stage
COPY --from=backend-builder /app/backend/tmp/main ./main

# Ensure the binary is executable
RUN chmod +x ./main

# Set environment variables
ENV CGO_ENABLED=1
ENV GIN_MODE=release

# Expose the application port
EXPOSE 8080

# Command to run the binary
CMD ["./main"]
