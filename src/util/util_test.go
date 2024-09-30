package util_test

import (
	"hudson-newey/2rm/src/util"
	"testing"
)

func TestInArrayTrue(t *testing.T) {
	expectedResult := true
	realizedResult := util.InArray([]int{1, 2, 3}, 1)

	if realizedResult != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, realizedResult)
	}
}

func TestInArrayFalse(t *testing.T) {
	expectedResult := false
	realizedResult := util.InArray([]int{1, 2, 3}, 4)

	if realizedResult != expectedResult {
		t.Fatalf("Expected %v but got %v", expectedResult, realizedResult)
	}
}
