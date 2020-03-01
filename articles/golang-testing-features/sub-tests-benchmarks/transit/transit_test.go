package transit

import (
	"testing"
)

func TestTransitRoutes(t *testing.T) {
	type TestData struct {
		firstStation string
		endStation string
		directTrain bool
	}
	tests := []TestData {
		{ firstStation: "Westminster", endStation: "Bond Street", directTrain: false },
		{ firstStation: "Temple", endStation: "Monument", directTrain: true },
		{ firstStation: "Barons Court", endStation: "Ladbroke Grove", directTrain: false },
		{ firstStation: "Colindale", endStation: "Belsize Park", directTrain: true },
	}
	tr := NewTransit()
	for _, test := range tests {
		direct := tr.Direct(test.firstStation, test.endStation)
		if direct != test.directTrain {
			t.Fatalf("expected: %v, got: %v", test.directTrain, direct)
		}
	}
}

func TestTransitRoutesSubtests(t *testing.T) {
	type TestData struct {
		line string
		firstStation string
		endStation string
		directTrain bool
	}
	tests := []TestData {
		{ line: "red,yellow", firstStation: "Westminster", endStation: "Bond Street", directTrain: false },
		{ line: "green,yellow", firstStation: "Temple", endStation: "Monument", directTrain: true },
		{ line: "blue,yellow", firstStation: "Barons Court", endStation: "Ladbroke Grove", directTrain: false },
		{ line: "black", firstStation: "Colindale", endStation: "Belsize Park", directTrain: true },
	}
	tr := NewTransit()
	for _, test := range tests {
		setup(t)
		t.Run(test.line, func(t *testing.T) {
			direct := tr.Direct(test.firstStation, test.endStation)
			if direct != test.directTrain {
				t.Fatalf("expected: %v, got: %v", test.directTrain, direct)
			}
		})
		teardown(t)
	}
}

func setup(t *testing.T) {
	// change global config
	t.Log("setup()")
	TransitConfig.Algorithm = BreadthFirst
}

func teardown(t *testing.T) {
	// cancel global config change
	t.Log("teardown()")
	TransitConfig.Algorithm = DepthFirst
}

func TestTransitRoutesSubtestsInParallel(t *testing.T) {
	
	type TestData struct {
		line string
		firstStation string
		endStation string
		directTrain bool
	}
	tests := []TestData {
		{ line: "red,yellow", firstStation: "Westminster", endStation: "Bond Street", directTrain: false },
		{ line: "green,yellow", firstStation: "Temple", endStation: "Monument", directTrain: true },
		{ line: "blue,yellow", firstStation: "Barons Court", endStation: "Ladbroke Grove", directTrain: false },
		{ line: "black", firstStation: "Colindale", endStation: "Belsize Park", directTrain: true },
	}
	tr := NewTransit()
	setup(t)
	t.Run("TestGroup", func(t *testing.T) {
		for _, test := range tests {
			test := test // capture the scope
			t.Run(test.line, func(t *testing.T) {
				t.Log(test.line, "start")
				t.Parallel()
				t.Log(test.line, "run")
				_ = tr.Direct(test.firstStation, test.endStation)
				t.Log(test.line, "end")
			})
		}
	})
	teardown(t)
	t.Log("TestTransitRoutesSubtestsInParallel end")
}

