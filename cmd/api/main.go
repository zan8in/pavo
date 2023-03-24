package main

import (
	"fmt"
	"strings"

	"github.com/zan8in/gologger"
	"github.com/zan8in/pavo/pkg/config"
	"github.com/zan8in/pavo/pkg/fofa"
	"github.com/zan8in/pavo/pkg/pavo"
	"github.com/zan8in/pavo/pkg/retryhttpclient"
)

func main() {

	options, err := pavo.NewOptions(pavo.Options{
		Query: []string{"server='thinkphp'"},
	})
	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}

	runner, err := pavo.NewRunner(options)
	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}

	runner.Run()

}

func main2() {

	gologger.Info().Msg("Pavo Running")

	retryhttpclient.Init(&retryhttpclient.Options{})

	c, err := config.NewConfig()
	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}
	fmt.Println(c.Fofa)

	fofa, err := fofa.New(&fofa.FofaOptions{Email: c.Fofa.Email, Key: c.Fofa.Key})
	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}

	fofa.SetSize(2)
	fofa.SetPage(2)
	fofa.SetFull(true)
	results, err := fofa.Query("server='thinkphp'")
	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}
	if len(results.Results) > 0 {
		for _, v := range results.Results {
			fmt.Println(strings.Join(v, ", "))
		}
	}

}
