package main

import (
	"net/http"
	"strconv"
	"time"
)

//Now is a stubbed accessor for time.Now
var Now = time.Now

//StddevResult contains standard deviation and associated input data
type StddevResult struct {
	Stddev float64 `json:"stddev,omitempty"`
	Data   []int   `json:"data,omitempty"`
}

func makeRandomAPIRequest(length int) []int {
	ret := []int{1, 2, 3, 4, 5}
	return ret

}

// askRandomAPI calls random.org API numReqs times, for length number of integers.
// Then, computes standard deviation for each response and standard deviation
// for numbers in all responses. Results are returned as a slice
// of StddevResult struct.
func askRandomAPI(numReqs, length int) []StddevResult {
	allNumbers := make([]int, 0)
	res := make([]StddevResult, 0)

	for i := 0; i < numReqs; i++ {
		numbers := makeRandomAPIRequest(5)
		allNumbers = append(allNumbers, numbers...)

		stddev := computeStdDev(numbers)
		res = append(res, StddevResult{Stddev: stddev, Data: numbers})
	}

	stddev := computeStdDev(allNumbers)
	res = append(res, StddevResult{Stddev: stddev, Data: allNumbers})

	return res
}

// getRandomMean serves GET requests to /random/mean endpoint.
// Mandatory query parameters: requests (int), length (int).
func getRandomMean(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	if !areParamsPresent(params, "requests", "length") {
		errorStruct := makeErrorStruct(http.StatusUnprocessableEntity,
			"Unprocessable Entity", "Missing parameters in query string",
			r.RequestURI)
		writeAsJsonWithStatus(errorStruct, errorStruct.Status, w)
		return
	}

	requests, errRequests := strconv.Atoi(params["requests"][0])
	length, errLength := strconv.Atoi(params["length"][0])

	if errRequests != nil || errLength != nil {
		errorStruct := makeErrorStruct(http.StatusUnprocessableEntity,
			"Unprocessable Entity", "Unable to parse parameters in query string",
			r.RequestURI)
		writeAsJsonWithStatus(errorStruct, errorStruct.Status, w)
		return
	}

	res := askRandomAPI(requests, length)
	writeAsJsonWithStatus(res, http.StatusOK, w)
}

func main() {
	http.HandleFunc("/random/mean", getRandomMean)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
