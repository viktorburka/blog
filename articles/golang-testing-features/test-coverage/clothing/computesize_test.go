package computesize

import "testing"

func TestComputeSize(t *testing.T) {
	type testData struct {
		chestSize int
		expected string
	}
	tests := []testData{
		{ chestSize: 38, expected: "S" },
		{ chestSize: 40, expected: "M" },
		{ chestSize: 40, expected: "M" },
		{ chestSize: 40, expected: "M" },
		{ chestSize: 40, expected: "M" },
		{ chestSize: 42, expected: "L" },
		{ chestSize: 42, expected: "L" },
	}
	for _, test := range tests {
		result := ComputeJacketSize(test.chestSize)
		if result != test.expected {
			t.Fatalf("expected size %v but got %v", test.expected, result)
		}
	}
}
