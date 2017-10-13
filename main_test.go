package main

import (
	"os"
	"testing"
)

// TestMain - test to drive external testing coverage
func TestMain(t *testing.T) {
	os.Args = os.Args[2:]
	main()
}
