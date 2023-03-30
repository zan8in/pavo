package main

import (
	"fmt"
	"strings"

	"github.com/zan8in/gologger"
	"github.com/zan8in/pavo/pkg/runner"
)

func main() {
	options, err := runner.ParseOptions()
	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}

	r, err := runner.NewRunner(options)
	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}

	if err := r.Run(); err != nil {
		gologger.Fatal().Msg(err.Error())
	}

	rs := r.Result.GetResult()
	for s := range rs {
		fmt.Println(strings.Join(s, ","))
	}

	runner.WriteOutput(r.Result)

}
