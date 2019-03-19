package main

import (
	"sync"
	"syscall/js"
)

type JsFuncTable struct {
	JsFunc func(int, string) (int, string)
}

var jsFuncTable *JsFuncTable

func goFunc(i int, s string) (int, string) {
	i, s = jsFuncTable.JsFunc(i, s)
	return i + 2, s + "c"
}

func main() {
	jsFuncs := js.Global().Get("jsFuncs")
	jsFuncTable = &JsFuncTable{
		JsFunc: func(i int, s string) (int, string) {
			res := jsFuncs.Get("jsFunc").Invoke(i, s)
			return res.Get("i").Int(), res.Get("s").String()
		},
	}

	goFuncs := js.Global().Get("goFuncs")
	goFuncs.Set("goFunc", js.NewCallback(func(args []js.Value) {
		i, s := goFunc(args[0].Int(), args[1].String())
		ret := args[2]
		ret.Set("i", i)
		ret.Set("s", s)
	}))

	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
