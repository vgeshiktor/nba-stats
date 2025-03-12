package integration_test

import (
	"fmt"
	"os"
	"testing"
)

func SetWorkDir() {
	projectRoot := os.Getenv("PROJECT_ROOT")
	if projectRoot == "" {
		projectRoot = "."
	}
	_ = os.Chdir(projectRoot)
}

func setup() {
	fmt.Println("Setting up before running tests...")
	// Add any necessary initialization logic here (e.g., database connection, config loading, etc.)
	SetWorkDir()
}

func teardown() {
	fmt.Println("Cleaning up after tests...")
	// Perform cleanup tasks here
}

func TestMain(m *testing.M) {
	// Run setup before any tests
	setup()

	// Run all tests
	exitCode := m.Run()

	// Run teardown after tests
	teardown()

	// Exit with the test run's result
	os.Exit(exitCode)
}