### GoLang Testing: Sub-tests and Sub-benchmarks

While writing unit tests in GoLang, I found few usefull features that may make your unit tests more structured and easier to manage. Here I would like to share with you the following: setup and teardown, control over test parallelism and finer control of table driven benchmarks and test cases.

#### Table driven tests and benchmarks

In GoLang table driven tests can be implemented via declaring an array of test cases first and then executing them in a loop. A typical example would look like:

```golang
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
```

There are few disadvantages of this approach. The first one is that its not possible to selectively run test cases out of the box and contol parallel execution unless unless writing some additinal code to handle that. However there is an easier way to do that if invoke some capabilities that GoLang testing package provides.

Another problem that stands out is when implementing test benchmarks in this fassion, the benchmark will measure the test suite as a whole rather than measuring the execution of each seprarate test case. A typical workaround is to implement each test case as a separate benchmark test case which calls some shared function with different params.

#### Fine control of table driven benchmarks and test cases

The above problems can be solved by creating subtests via calling `t.Run` function. It takes a name string parameter as a first argument and a function to run as a second one. When a unit test is run this way, it will logically seprate each of `t.Run` invocations as a seprate sub-test or sub-benchmark.

```golang
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
		t.Run(test.line, func(t *testing.T) {
			direct := tr.Direct(test.firstStation, test.endStation)
			if direct != test.directTrain {
				t.Fatalf("expected: %v, got: %v", test.directTrain, direct)
			}
		})
	}
}
```
With that setup, we can now run subtests that only validate the green line:

```bash
$ go test -v -run=TestTransitRoutesSubtests/green ./...
=== RUN   TestTransitRoutesSubtests
=== RUN   TestTransitRoutesSubtests/green,yellow
--- PASS: TestTransitRoutesSubtests (0.00s)
    --- PASS: TestTransitRoutesSubtests/green,yellow (0.00s)
PASS
ok  	_/transit	0.005s
```

#### Setup and teardown

For sequential tests, such as the above, setup and teardown code is pretty straightfoward. To do that, we just declare setup and teardown functions and call them before and after the place where the test code is run. Since `t.Run` method is guaraneed to return only after the test function is finished, setup and teardown functions will always be called sequentially. This comes handy when we need to maintain same global state since each new unit test should be invoked with the original state and not altered by the previous test.

```golang
func TestTransitRoutesSubtests(t *testing.T) {
    //...
	tr := NewTransit()
	setup()
	for _, test := range tests {
		t.Run(test.line, func(t *testing.T) {
			direct := tr.Direct(test.firstStation, test.endStation)
			if direct != test.directTrain {
				t.Fatalf("expected: %v, got: %v", test.directTrain, direct)
			}
		})
	}
	teardown()
}

func setup() {
	// change global config
	TransitConfig.Algorithm = BreadthFirst
}

func teardown() {
	// cancel global config change
	TransitConfig.Algorithm = DepthFirst
}
```


#### Control over test parallelism

Lets now look at the sutuation where we need a parallel execution of our unit test subtests. As a reminder, in the previous example all subtests were executed sequentially, which means that any next subtest wouldn't start until the previous one is completed. If we want to test for a correct concurrency execution in that scenario, say with `-race` flag to detect race couditions that wouldn't be possible since `-race` option requires a concurrent code to be running to be to detect any anomalies.

Since in the previous example a sequential execution was guaranteed, the code would never run concurrently. One way to fix that would be to spawn multiple goroutines but then we would need to control syncronization ourselves adding additional code which would mix with the test code itself reducing the readability of the test. However, there is an easier way to do that and its to invoke a `t.Parallel` method.

```golang
func TestTransitRoutesSubtestsInParallel(t *testing.T) {
	//...
	tr := NewTransit()
	for _, test := range tests {
		test := test // capture the scope
		t.Run(test.line, func(t *testing.T) {
			t.Parallel()
			direct := tr.Direct(test.firstStation, test.endStation)
			if direct != test.directTrain {
				t.Errorf("expected: %v, got: %v", test.directTrain, direct)
			}
		})
	}
}
```

This time we run all our subtest in parallel rather than sequentially and the unit test itself won't finish until all the subtests are done. Note that in this case the order of subtest execution is not guaranteed. What is guaranteed is that all subtests will finish when the unit test returns and that all of them running concurrently won't interfere with subtests of other unit tests that run concurrently.

Lets now make sure that `setup` and `teardown` functions will work if applied to this code. We can't use them in the same way we did in the previous case because even though the completion of the subtests is guaranteed, the serialization of statements within the unit tests itself and the end of subtests is not. In other words the unit test can reach the end before all the subtest will be completed. To serialize that part, we can call `t.Run` method twice.

```golang
func TestTransitRoutesSubtestsInParallel(t *testing.T) {
	//...
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
}
```
We add scope capturing code to ensure the `test` value doesn't get overwritten and the vairable doesn't get involved in a race condition. Lets run it now and make sure we see our `setup` and `teardown` being invoked properly. 

```bash
$ go test -v -run=TestTransitRoutesSubtestsInParallel/TestGroup/green ./...
...
--- PASS: TestTransitRoutesSubtestsInParallel (0.00s)
    transit_test.go:56: setup()
    --- PASS: TestTransitRoutesSubtestsInParallel/TestGroup (0.00s)
        --- PASS: TestTransitRoutesSubtestsInParallel/TestGroup/green,yellow (0.00s)
            transit_test.go:86: green,yellow start
            transit_test.go:88: green,yellow run
            transit_test.go:90: green,yellow end
    transit_test.go:62: teardown()
PASS
ok  	_/transit	0.006s
```
