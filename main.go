package main

import (
	"encoding/json"
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

	_, requestsParamPresent := params["requests"]
	_, lengthParamPresent := params["length"]

	if !requestsParamPresent || !lengthParamPresent {
		status := http.StatusUnprocessableEntity
		w.WriteHeader(status)

		errJSON, err := json.Marshal(makeErrorStruct(status, "Unprocessable Entity",
			"Missing parameters in query string", r.RequestURI))
		if err == nil {
			w.Write(errJSON)
		}

		return
	}

	reqNum, errReqNum := strconv.Atoi(params["requests"][0])
	length, errLength := strconv.Atoi(params["length"][0])

	if errReqNum != nil || errLength != nil {
		status := http.StatusUnprocessableEntity
		w.WriteHeader(status)

		errJSON, err := json.Marshal(makeErrorStruct(status, "Unprocessable Entity",
			"Unable to parse parameters in query string", r.RequestURI))
		if err == nil {
			w.Write(errJSON)
		}

		return
	}

	res := askRandomAPI(reqNum, length)
	jsonres, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonres)
}

func main() {
	http.HandleFunc("/random/mean", getRandomMean)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
