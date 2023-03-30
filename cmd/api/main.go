package main

import (
	"fmt"

	"github.com/zan8in/pavo"
)

func main() {
	r, err := pavo.QuerySubDomain("example.com")
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, v := range r {
		fmt.Println(v)
	}
}
