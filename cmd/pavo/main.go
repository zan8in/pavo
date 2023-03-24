package main

import (
	"fmt"

	"github.com/zan8in/gologger"
	"github.com/zan8in/pavo/pkg/pavo"
)

func main() {
	options, err := pavo.ParseOptions()
	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}

	runner, err := pavo.NewRunner(options)
	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}

	if err := runner.Run(); err != nil {
		gologger.Fatal().Msg(err.Error())
	}

	rs := runner.Result.GetResult()
	for s := range rs {
		fmt.Println(s)
	}

}
