package main

import "math"

func reverse(x int) int {
	if int(math.Pow(2, 31))/10-1 < x || x < int(math.Pow(-2, 31))/10 {
		return 0
	}
	result := []int{}
	var res int = 0
	for {
		j := x % 10
		x = x / 10
		result = append(result, j)
		if x == 0 {
			break
		}
	}
	if len(result) > 10 {
		return 0
	}
	for i := 0; i < len(result); i++ {
		res = res*10 + result[i]
	}
	return res
}
