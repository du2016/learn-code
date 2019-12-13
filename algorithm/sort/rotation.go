//手摇算法
package main

import "log"

func main() {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	rotation(a, 4, 8)
	log.Println(a)
}

func reverse(a []int, n int) {
	if n < 2 {
		return
	}
	for i := 0; i < n; i++ {
		n--
		a[i], a[n] = a[n], a[i]
	}
}

func rotation(a []int, m, n int) {
	reverse(a, m)
	reverse(a, n-m)
	reverse(a, n)
}
