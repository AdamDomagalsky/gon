package api

// go test ./... -v -run TestCounterHandlers/CounterScopedValue -race;
import (
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"sync"
	"testing"
)

func TestCounterHandlers(t *testing.T) {
	testCases := []struct {
		name       string
		iterations int
		endpoint   string
	}{
		{
			name:       "Basic",
			iterations: 99999,
			endpoint:   "/api/v1/counter/bad",
		},
		{
			name:       "Scoped",
			iterations: 99999,
			endpoint:   "/api/v1/counter/scoped",
		},
		{
			name:       "Mutex",
			iterations: 99999,
			endpoint:   "/api/v1/counter/mutex",
		},
		{
			name:       "Atomic",
			iterations: 99999,
			endpoint:   "/api/v1/counter/atomic",
		},
		{
			name:       "Semaphore",
			iterations: 99999,
			endpoint:   "/api/v1/counter/semaphore",
		},
		{
			name:       "Channel",
			iterations: 999999,
			endpoint:   "/api/v1/counter/channel",
		},
		// Add more test cases here for other endpoints if needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			results := coreIteration(t, tc.iterations, tc.endpoint)
			testForDuplicateCounts(t, results)
		})
	}
}

func coreIteration(t *testing.T, iterations int, endpoint string) []int {
	// Create a new server instance
	server := NewServer(&Config{Port: "8080"})

	// Create a test request to the /CounterMutex endpoint
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}
	results := []int{}
	resultsMutex := sync.Mutex{}
	var wg sync.WaitGroup
	for i := range iterations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()
			// Serve the request
			server.router.ServeHTTP(rr, req)
			// Check the status code
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("%d Handler returned wrong status code: got %v want %v", i, status, http.StatusOK)
			}
			n, _ := strconv.Atoi(rr.Body.String())
			resultsMutex.Lock()
			results = append(results, n)
			resultsMutex.Unlock()
		}()
	}
	wg.Wait()
	if len(results) != iterations {
		t.Errorf("Expected %d results, got %d", iterations, len(results))
	}
	return results
}

func testForDuplicateCounts(t *testing.T, results []int) {
	// isVerbose := testing.Verbose()

	// Create a map to count occurrences of each number
	countMap := make(map[int]int)
	for _, num := range results {
		countMap[num]++
	}
	// Create a sorted slice of keys
	keys := make([]int, 0, len(countMap))
	for k := range countMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	foundDuplicates := false
	for _, num := range keys {
		count := countMap[num]
		if count > 1 {
			// if isVerbose {
			// 	t.Logf("Number %d appears %d times\n", num, count)
			// }
			foundDuplicates = true
		}
	}
	if foundDuplicates {
		t.Errorf("Found duplicates in the results. To see the duplicates, run the test with the -v flag")
	}
}
