### GoLang Testing: Race Conditions

It is needless to say how important the unit testing is in early detection of potential problems with our code. And GoLang team seeing value in unit tests provides nice infrastructure for writing unit tests out of the box with GoLang SDK. The capabilities are extensive and go way beyond just validating the results in a traditional unit test application where we run a unit test, then compare the results with expected values and either stop or go to the next test. Here, I will try to talk about some of them that I find quite useful.

With goroutines being a part of the language itself and such a commonly used feature in Go applications, being able to validate our code for potential race conditions is a really good thing to have. The `go test` command comes with a `-race` option which instruments the code during compilation to allow it detect race condition situations. The disadvantage and the reason it is disabled by default is because the instrumentation slows down the code execution. It therefore leads to a slower unit test execution which may be an issue for projects with large code and tests base. An optimum solution would probably be to have it enabled for specific tests where the test execution time is negligible comparing to the benefit it brings.

Let’s look at an example where we have a member variable of some struct which can be accessed from multiple methods and doesn't have a syncronization primitive protecting it from a concurrent access.

```golang
type NodeAccessor struct {
	counter map[string]int
}

func NewNodeAccessor() *NodeAccessor {
	return &NodeAccessor{
		counter: make(map[string]int),
	}
}

func (na *NodeAccessor) UpdateNode(name string) {
	if _, ok := na.counter[name]; !ok {
		na.counter[name] = 0
	} else {
		na.counter[name]++
	}
}
```

If we run the following test case for `UpdateNode` method, it will pass because the access to the methods is non-concurrent and therefore there is no race condition situation.

```golang
func TestUpdateNodeNonConcurrent(t *testing.T) {
	type TestData struct {
		nodeName string
		expectedNodeCount int
	}
	tests := []TestData{
		{nodeName: "/stream/1", expectedNodeCount: 0},
		{nodeName: "/stream/2", expectedNodeCount: 0},
		{nodeName: "/stream/1", expectedNodeCount: 1},
		{nodeName: "/stream/1", expectedNodeCount: 2},
	}
	accessor := NewNodeAccessor()
	for _, test := range tests {
		accessor.UpdateNode(test.nodeName)
		nodeCount, err := accessor.NodeCount(test.nodeName)
		if err != nil {
			t.Fatalf("node %v not found", test.nodeName)
		}
		if nodeCount != test.expectedNodeCount {
			t.Fatalf("incorrect node count for node %v: expected %v but got %v",
			 test.nodeName, test.expectedNodeCount, nodeCount)
		}
	}
}
```

```bash
$ go test -race
PASS
ok  	_/golang-testing-features	1.010s
```

We need to change the test to introduce concurrent access to `UpdateNode` method. It can be done in several ways. One of them would be to call the method in a seprate goroutine each time for every test. We will then need to make sure the goroutines complete and do not exit prematurely. Another, more idomatic way from the test perspective is to use `Parallel` method which uses a technique introduced in the previous section.

```golang
for idx, test := range tests {
	t.Run(fmt.Sprintf("test#%v",idx), func(t *testing.T) {
		test := test // capture range variable
    	t.Parallel()
		accessor.UpdateNode(test.nodeName)
		nodeCount, err := accessor.NodeCount(test.nodeName)
		if err != nil {
			t.Fatalf("node %v not found", test.nodeName)
		}
		...
    })
}
```

This time the test will fail since data race will be detected.

```bash
$ go test -race
==================
WARNING: DATA RACE
Read at 0x00c0000922a0 by goroutine 10:
  runtime.mapaccess2_faststr()
      /usr/local/go/src/runtime/map_faststr.go:101 +0x0
  _/golang-testing-features.TestUpdateNodeConcurrent.func1()
      /golang-testing-features/nodeaccessor.go:18 +0x101
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:827 +0x162
...
==================
--- FAIL: TestUpdateNodeConcurrent (0.00s)
    --- FAIL: TestUpdateNodeConcurrent/test#0 (0.00s)
    --- FAIL: TestUpdateNodeConcurrent/test#1 (0.00s)
        testing.go:771: race detected during execution of test
    --- FAIL: TestUpdateNodeConcurrent/test#2 (0.00s)
        testing.go:771: race detected during execution of test
    --- FAIL: TestUpdateNodeConcurrent/test#3 (0.00s)
        nodeaccessor_test.go:55: incorrect node count for node /stream/1: expected 2
        but got 1
FAIL
exit status 1
FAIL	_/golang-testing-features	0.011s
```

We can fix this problem by synchronizing access to the internal `counter` map. The most straightforward way to do that is to add a mutex variable that will control serialized access to the structure. If same `UpdateNode` method is called from different goroutines, the mutex `Lock` method will lock the rest of them while the first goroutine which acquired access to the mutex runs the method. Upon unlocking, the next awaiting goroutine will acquire the mutex and perform the operation on `UpdateNode` method. We also use `defer` method as a handy way to do not forget to unlock the mutex.

```golang
type NodeAccessor struct {
    counter map[string]int
    mux sync.Mutex
}

func (na *NodeAccessor) UpdateNode(name string) {
	na.mux.Lock()
	defer na.mux.Unlock()

	if _, ok := na.counter[name]; !ok {
		na.counter[name] = 0
	} else {
		na.counter[name]++
	}
}


func (na *NodeAccessor) NodeCount(name string) (int, error) {
	na.mux.Lock()
	defer na.mux.Unlock()
	...
}

```

```bash
$ go test -race
PASS
ok  	_/golang-testing-features	1.012s
```
