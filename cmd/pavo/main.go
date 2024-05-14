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
	k := 1
	for s := range rs {
		// fmt.Printf("s[0]=%s,s[1]=%s,s[2]=%s,s[3]=%s,s[4]=%s\n", s[0], s[1], s[2], s[3], s[4])
		fmt.Println(k, "=========", strings.Join(s, ","))
		k++
	}

	runner.WriteOutput(r.Result)

}
