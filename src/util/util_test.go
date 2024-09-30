package util

import "testing"

func TestInArrayTrue(t *testing.T) {
	expectedResult := true
	realizedResult := InArray([]int{1, 2, 3}, 1)

	if realizedResult != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, realizedResult)
	}
}

func TestInArrayFalse(t *testing.T) {
	expectedResult := false
	realizedResult := InArray([]int{1, 2, 3}, 4)

	if realizedResult != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, realizedResult)
	}
}
