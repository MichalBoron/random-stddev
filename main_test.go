package main

import (
	"math"
	"testing"
	"time"
)

func SetNow() {
	Now = func() time.Time {
		return time.Date(1980, 6, 1, 11, 0, 0, 0, time.UTC)
	}

}

func TestMakeErrorStruct(t *testing.T) {
	SetNow()
	got := makeErrorStruct(500, "errorMsg", "message", "path")
	want := ErrorStruct{
		Timestamp: Now().String(),
		Status:    500,
		ErrorMsg:  "errorMsg",
		Message:   "message",
		Path:      "path",
	}

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}

func TestComputeStdDev(t *testing.T) {
	var testCases = []struct {
		nums []int
		want float64
	}{
		{[]int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5}, 1.4142},
		{[]int{1, 1, 1, 1}, 0.0},
		{[]int{1, 1, 3, 3}, 1.0},
		{[]int{3, 5, 9, 1, 8, 6, 58, 9, 4, 10}, 15.811704},
	}

	for _, test := range testCases {
		got := computeStdDev(test.nums)
		want := test.want

		if math.Abs(got-want) > 0.0001 {
			t.Errorf("got %f, wanted %f", got, want)
		}
	}

}
