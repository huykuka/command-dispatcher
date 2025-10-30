package hashing

import (
	"errors"
	"os"
	"testing"
)

func TestHashingServiceSingleton(t *testing.T) {
	// Set up the environment variable for the password key
	os.Setenv("HASH_PASSWORD_KEY", "supersecretkey")

	// Call NewHashingService multiple times
	service1 := NewHashingService()
	service2 := NewHashingService()

	// Assert that both instances point to the same object (singleton behavior)
	if service1 != service2 {
		t.Errorf("Expected singleton instance, but got different instances")
	}
}

func TestHashPasswordAndComparePasswords_TableDriven(t *testing.T) {
	// Set up the environment variable for the password key
	os.Setenv("HASH_PASSWORD_KEY", "supersecretkey")

	// Get the singleton instance of HashingService
	service := NewHashingService()

	// Define a test table with "should" in the test names
	tests := []struct {
		name           string
		password       string
		hashedPassword string
		expectedError  error
	}{
		{
			name:          "should successfully compare valid password",
			password:      "mysecretpassword",
			expectedError: nil,
		},
		{
			name:           "should fail to compare invalid password",
			password:       "wrongpassword",
			hashedPassword: "fakehash", // A fake hash to simulate a failed comparison
			expectedError:  errors.New("invalid password"),
		},
	}

	// Iterate over the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// If a hashedPassword is not predefined, generate it
			if tt.hashedPassword == "" {
				tt.hashedPassword = service.HashPassword(tt.password)
			}

			// Compare the password with the hashed password
			err := service.ComparePasswords(tt.hashedPassword, tt.password)

			// Check if the error matches the expected error
			if err != nil && tt.expectedError == nil {
				t.Errorf("Expected no error, but got: %v", err)
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("Expected error '%v', but got none", tt.expectedError)
			}
			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("Expected error '%v', but got '%v'", tt.expectedError, err)
			}
		})
	}
}
