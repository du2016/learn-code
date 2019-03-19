package main

import (
	"fmt"
	"unsafe"
)

func main() {
	//var a int64 = 123
	//var b float64 = float64(a)
	//fmt.Println(&a)
	//fmt.Println(&b)

	//var a int64 = 123
	//pn := &a
	//pf := (*float64)(unsafe.Pointer(pn))
	//*pf = 3.5
	//fmt.Println(a)
	//fmt.Println(*pn)
	//fmt.Println(*pf)

	a := [4]int{0, 1, 2, 3}
	p1 := unsafe.Pointer(&a[1])
	p3 := unsafe.Pointer(uintptr(p1) + 2*unsafe.Sizeof(a[3]))
	*(*int)(p3) = 4
	fmt.Println(a)
}
