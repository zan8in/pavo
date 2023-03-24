package pavo

import (
	"fmt"
	"strings"
	"time"

	"github.com/remeh/sizedwaitgroup"
	"github.com/zan8in/gologger"
	"github.com/zan8in/pavo/pkg/config"
	"github.com/zan8in/pavo/pkg/fofa"
	"github.com/zan8in/pavo/pkg/result"
	"github.com/zan8in/pavo/pkg/retryhttpclient"
)

type (
	Runner struct {
		options *Options
		config  *config.Config
		fofa    *fofa.FofaOptions
		ticker  *time.Ticker
		wgscan  sizedwaitgroup.SizedWaitGroup
		Result  *result.Result
	}
)

func NewRunner(options *Options) (*Runner, error) {
	r := &Runner{
		options: options,
		Result:  result.NewResult(),
	}

	config, err := config.NewConfig()
	if err != nil {
		return r, err
	}
	r.config = config

	retryhttpclient.Init(&retryhttpclient.Options{
		Timeout: DefaultTimeout,
		Retries: DefaultRetries,
	})

	if err := r.initPlatform(); err != nil {
		return r, err
	}

	r.ticker = time.NewTicker(time.Second / time.Duration(DefaultRateLimit))
	r.wgscan = sizedwaitgroup.New(DefaultRateLimit)

	return r, nil
}

func (r *Runner) Run() error {
	if r.config.IsFofa() && r.options.Platform == FofaPlatform {
		r.RunFofa()
	}
	if r.options.Platform == HunterPlatform {
		fmt.Println("this is hunter platform")
	}

	return nil
}

func (r *Runner) RunFofa() {
	r.Result.AddQuery(strings.Join(r.options.Query, ","))

	for _, q := range r.options.Query {
		r.wgscan.Add()
		go func(q string) {
			defer r.wgscan.Done()
			<-r.ticker.C

			if r.options.Count > DefaultQueryCount {
				page := 1
				for {
					n := r.options.Count - page*DefaultQueryCount
					// fmt.Printf("%d - %d = %d..............", r.options.Count, page*DefaultQueryCount, n)
					if n >= 0 {
						// fmt.Printf("page=%d&size=%d\n", page, DefaultQueryCount)
						r.fofa.ReSet()
						r.fofa.SetPage(page)
						r.fofa.SetSize(DefaultQueryCount)
						results, err := r.fofa.Query(q)
						if err != nil {
							gologger.Fatal().Msg(err.Error())
						}
						if len(results.Results) > 0 {
							r.Result.AddResult(results.Results)
						}
					} else {
						last := DefaultQueryCount - (page*DefaultQueryCount - r.options.Count)
						if last <= 0 {
							break
						}
						// fmt.Printf("page=%d&size=%d\n", page, last)
						r.fofa.ReSet()
						r.fofa.SetPage(page)
						r.fofa.SetSize(last)
						results, err := r.fofa.Query(q)
						if err != nil {
							gologger.Fatal().Msg(err.Error())
						}
						if len(results.Results) > 0 {
							r.Result.AddResult(results.Results)
						}
						break
					}
					page++
					time.Sleep(100 * time.Millisecond)
				}
			} else {
				r.fofa.ReSet()
				r.fofa.SetSize(r.options.Count)
				results, err := r.fofa.Query(q)
				if err != nil {
					gologger.Fatal().Msg(err.Error())
				}
				if len(results.Results) > 0 {
					r.Result.AddResult(results.Results)
				}
			}

		}(q)
	}
	r.wgscan.Wait()
}

func (r *Runner) initPlatform() (err error) {
	if !r.config.IsFofa() {
		return fmt.Errorf("missing fofa email and key")
	}

	fofaEmail, fofaKey := r.config.Fofa.Email, r.config.Fofa.Key
	if fofa, err := fofa.New(&fofa.FofaOptions{Email: fofaEmail, Key: fofaKey}); err == nil {
		r.fofa = fofa
		r.fofa.SetSize(r.options.Count)
		return nil
	}

	return err
}
