package main

import (
	"fmt"
	"testing"
)

type mockData string

var mockdata mockData

func TestAddWithTestingPackage(t *testing.T) {
	prepare()
	result := TestedFunction()
	fmt.Println(result)
	if result != 3 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", result, 3)
	}
}

func prepare() {
	mockdata = mockData("mock")
}

func (m mockData) multiply(a, b int) int {
	return a * b * b
}
