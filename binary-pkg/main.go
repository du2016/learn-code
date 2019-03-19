package main

import (
	"fmt"
	"github.com/du2016/learn-code/binary-pkg/x"
)

func main() {
	p := x.Person{
		Name: "test",
	}
	fmt.Println(p.Name)
	//fmt.Printf("result=%d\n", x.Add(1, 2))
}
