package main

import (
	"fmt"
	"log"
)

func main() {
	a := []int{2, 1, 3, 4, 6, 5, 7, 10}
	//shellsort(a)
	aux := make([]int, len(a))
	mergesort(a, aux, 0, 3, len(a))
	log.Println(a)
}

func bubblesort(a []int) {
	for i := 0; i < len(a); i++ {
		for j := i + 1; j < len(a); j++ {
			if a[j] < a[i] {
				a[i], a[j] = a[j], a[i]
			}
		}
	}
}

func selectsort(a []int) {
	for i := 0; i < len(a); i++ {
		min := i
		for j := i + 1; j < len(a); j++ {
			if a[j] < a[i] {
				min = j
			}
		}
		if min != i {
			a[i], a[min] = a[min], a[i]
		}
	}
}

func insertsort(a []int) {
	for i := 1; i < len(a); i++ {
		if a[i] < a[i-1] {
			j := i - 1
			temp := a[i]
			for j >= 0 && a[j] > temp {
				a[j+1] = a[j]
				j--
			}
			a[j+1] = temp
		}
	}
}

func binarysearch(a []int, k int) int {
	low := 0
	high := len(a)
	for low < high {
		mid := (low + high) / 2
		if a[mid] > k {
			high = mid - 1
		} else if a[mid] < k {
			low = mid + 1
		} else {
			return mid
		}
	}
	return -1
}

func shellsort(a []int) {
	N := len(a)
	h := 1
	for h < N/3 {
		h = 3*h + 1
	}
	for h >= 1 {
		for i := h; i < N; i++ {
			for j := i; j >= h && a[j] < a[j-h]; j -= h {
				a[j], a[j-h] = a[j-h], a[j]
			}
		}
		h = h / 3
	}
}

func Quick2Sort(values []int) {
	if len(values) <= 1 {
		return
	}
	mid, i := values[0], 1
	head, tail := 0, len(values)-1
	for head < tail {
		fmt.Println(values)
		if values[i] > mid {
			values[i], values[tail] = values[tail], values[i]
			tail--
		} else {
			values[i], values[head] = values[head], values[i]
			head++
			i++
		}
	}
	values[head] = mid
	Quick2Sort(values[:head])
	Quick2Sort(values[head+1:])
}

func mergesort(s, aux []int, lo, mid, hi int) {
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
