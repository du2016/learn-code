package main

func twoSum(nums []int, target int) []int {
	length := len(nums)
	for i := 0; i <= length-1; i++ {
		for j := i + 1; j <= length-1; j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return []int{}
}
