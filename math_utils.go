package main

import "math"

func computeStdDev(nums []int) float64 {
	var res float64 = 0.0
	count := len(nums)
	sum := 0

	for _, n := range nums {
		sum += n
	}

	var mean float64 = float64(sum) / float64(count)

	for _, n := range nums {
		partialRes := (float64(n) - mean) * (float64(n) - mean)
		res += partialRes
	}
	res = res / float64(count)
	res = math.Sqrt(res)

	return res
}
