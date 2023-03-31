package main

import (
	"fmt"

	"github.com/zan8in/pavo"
)

func main() {
	r, err := pavo.QueryIPPort("60.10.113.44", 100)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, v := range r.Result {
		fmt.Println(v)
	}

	// r, err := pavo.QuerySubDomain("example.com", 1)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// for _, v := range r {
	// 	fmt.Println(v)
	// }
}
