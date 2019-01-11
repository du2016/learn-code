//原地归并
package main

import "fmt"

//自顶向下归并排序
func MergeSortUpToDown(s []int) {
	aux := make([]int, len(s)) //辅助切片
	mergeSortUpToDown(s, aux, 0, len(s)-1)
}

//自底向上归并排序
func MergeSortDownToUp(s []int) {
	aux := make([]int, len(s)) //辅助切片
	n := len(s)
	for sz := 1; sz < n; sz *= 2 {
		for lo := 0; lo < n-sz; lo += 2 * sz {
			merge(s, aux, lo, lo+sz-1, min(lo+2*sz-1, n-1))
		}
	}
}

func mergeSortUpToDown(s, aux []int, lo, hi int) {
	if lo >= hi {
		return
	}
	mid := (lo + hi) >> 1
	mergeSortUpToDown(s, aux, lo, mid)
	mergeSortUpToDown(s, aux, mid+1, hi)
	merge(s, aux, lo, mid, hi)
}

//归并操作
func merge(s, aux []int, lo, mid, hi int) {
	for k := lo; k <= hi; k++ {
		aux[k] = s[k]
	}
	i, j := lo, mid+1
	for k := lo; k <= hi; k++ {
		if i > mid {
			s[k] = aux[j]
			j++
		} else if j > hi {
			s[k] = aux[i]
			i++
		} else if aux[j] < aux[i] {
			s[k] = aux[j]
			j++
		} else {
			s[k] = aux[i]
			i++
		}
	}
}

func min(i, j int) int {
	if j < i {
		return j
	}
	return i
}
func main() {
	s := []int{9, 0, 6, 5, 8, 2, 1, 7, 4, 3}
	fmt.Println(s)
	//MergeSortUpToDown(s)
	MergeSortDownToUp(s)
	fmt.Println(s)

}
