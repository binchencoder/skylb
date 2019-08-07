package main

import (
	"testing"
)

// TestFlags checks if there are flag duplication, which will cause application
// panic during startup, and is not detected during compile stage.
//
// To verify this test works, you can add the following line to main.go:
// v = flag.String("v", "", "A flag intentionally added to be duplicate with glog.")
func TestFlags(t *testing.T) {
	_ = main
}
