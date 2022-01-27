package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//Now is a stubbed accessor for time.Now
var Now = time.Now

//StddevResult contains standard deviation and associated input data
type StddevResult struct {
	Stddev float64 `json:"stddev,omitempty"`
	Data   []int   `json:"data,omitempty"`
}

//getIntsFromRandomAPI sends one GET request to random.org API
// asking for length number of integers in range 0 to 100.
func getIntsFromRandomAPI(length int) ([]int, error) {
	ret := make([]int, length)
	apiUrlTemplate := "https://www.random.org/integers/?num=%d&min=0&max=100&col=1&base=10&format=plain"
	apiUrl := fmt.Sprintf(apiUrlTemplate, length)

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(),
		time.Second*5)

	defer cancel()

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, errors.New("error in creating request")
	}

	response, err := http.DefaultClient.Do(req.WithContext(ctxWithTimeout))

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("error in get request to random.org API")
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("error reading response from random.org API")
	}

	numsString := strings.Split(string(respBody), "\n")
	for i, ns := range numsString {
		if i < length {
			n, err := strconv.Atoi(ns)
			if err != nil {
				return nil, errors.New("error parsing response from random.org API")
			}
			ret[i] = n
		}
	}
	return ret, nil
}

// askRandomAPI calls random.org API numReqs times, for length number of integers.
// Then, computes standard deviation for each response and standard deviation
// for numbers in all responses. Results are returned as a slice
// of StddevResult struct.
func askRandomAPI(numReqs, length int) ([]StddevResult, error) {
	allNumbers := make([]int, 0)
	res := make([]StddevResult, 0)

	for i := 0; i < numReqs; i++ {
		numbers, err := getIntsFromRandomAPI(5)
		if err != nil {
			return nil, err
		}
		allNumbers = append(allNumbers, numbers...)

		stddev := computeStdDev(numbers)
		res = append(res, StddevResult{Stddev: stddev, Data: numbers})
	}

	stddev := computeStdDev(allNumbers)
	res = append(res, StddevResult{Stddev: stddev, Data: allNumbers})

	return res, nil
}

// getRandomStddev serves GET requests to ednpoints:
//  /random/stddev and /random/mean
// Mandatory query parameters: requests (int), length (int).
func getRandomStddev(w http.ResponseWriter, r *http.Request) {
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

	res, err := askRandomAPI(requests, length)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	writeAsJsonWithStatus(res, http.StatusOK, w)
}

func main() {
	http.HandleFunc("/random/mean", getRandomStddev)
	http.HandleFunc("/random/stddev", getRandomStddev)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
