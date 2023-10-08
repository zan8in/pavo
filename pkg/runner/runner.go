package runner

import (
	"fmt"
	"strings"
	"time"

	"github.com/remeh/sizedwaitgroup"
	"github.com/zan8in/gologger"
	"github.com/zan8in/pavo/pkg/config"
	"github.com/zan8in/pavo/pkg/fofa"
	"github.com/zan8in/pavo/pkg/hunter"
	"github.com/zan8in/pavo/pkg/result"
	"github.com/zan8in/pavo/pkg/retryhttpclient"
)

type (
	Runner struct {
		options *Options
		config  *config.Config
		fofa    *fofa.FofaOptions
		hunter  *hunter.HunterOptions
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
	} else if r.config.IsHunter() && r.options.Platform == HunterPlatform {
		r.RunHunter()
	} else {
		return fmt.Errorf("no supported platform (fofa, Hunter)")
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
							gologger.Error().Msg(err.Error())
							return
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
							gologger.Error().Msg(err.Error())
							return
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
					gologger.Error().Msg(err.Error())
					return
				}
				if len(results.Results) > 0 {
					r.Result.AddResult(results.Results)
				}
			}

		}(q)
	}
	r.wgscan.Wait()
}

func (r *Runner) RunHunter() {
	r.Result.AddQuery(strings.Join(r.options.Query, ","))

	for _, q := range r.options.Query {
		r.wgscan.Add()
		go func(q string) {
			defer r.wgscan.Done()
			<-r.ticker.C

			rlist := []*hunter.HunterResultList{}

			if r.options.Count > DefaultQueryCount {
				page := 1
				for {
					n := r.options.Count - page*DefaultQueryCount
					// fmt.Printf("++++++%d - %d = %d..............", r.options.Count, page*DefaultQueryCount, n)
					if n >= 0 {
						// fmt.Printf("++++++page=%d&size=%d\n", page, DefaultQueryCount)
						r.hunter.ReSet()
						r.hunter.SetPage(page)
						r.hunter.SetSize(DefaultQueryCount)
						results, err := r.hunter.Query(q)
						if err != nil {
							gologger.Error().Msg(err.Error())
							return
						}
						if results != nil {
							rlist = append(rlist, results)
						}
						// if results != nil && len(results.Data.Arr) > 0 {
						// 	r.Result.AddResult(r.hunter.HunterResultList2Slice(results))
						// }
					} else {
						last := DefaultQueryCount - (page*DefaultQueryCount - r.options.Count)
						if last <= 0 {
							break
						}
						// fmt.Printf("--------page=%d&size=%d\n", page, last)
						r.hunter.ReSet()
						r.hunter.SetPage(page)
						r.hunter.SetSize(last)
						results, err := r.hunter.Query(q)
						if err != nil {
							gologger.Error().Msg(err.Error())
							return
						}
						if results != nil {
							rlist = append(rlist, results)
						}
						// if results != nil && len(results.Data.Arr) > 0 {
						// 	r.Result.AddResult(r.hunter.HunterResultList2Slice(results))
						// }
						break
					}
					page++
					time.Sleep(6000 * time.Millisecond)
				}
			} else {
				r.hunter.ReSet()
				r.hunter.SetPage(1)
				r.hunter.SetSize(r.options.Count)
				results, err := r.hunter.Query(q)
				if err != nil {
					gologger.Error().Msg(err.Error())
					return
				}
				if results != nil {
					rlist = append(rlist, results)
				}
				// if results != nil && len(results.Data.Arr) > 0 {
				// 	r.Result.AddResult(r.hunter.HunterResultList2Slice(results))
				// }
			}

			if len(rlist) > 0 {
				count := 0
				for _, rst := range rlist {
					count += len(rst.Data.Arr)
					r.Result.AddResult(r.hunter.HunterResultList2Slice(rst))
				}
				// fmt.Println("xxxxxxxxxxxxx:", count)
			}

		}(q)
	}
	r.wgscan.Wait()
}

func (r *Runner) initPlatform() (err error) {

	if r.options.Platform == HunterPlatform {
		if !r.config.IsHunter() {
			return fmt.Errorf("missing hunter api-key")
		}

		if r.config.IsHunter() {

			for _, key := range r.config.Hunter.ApiKey {
				if hunter, err := hunter.New(&hunter.HunterOptions{Key: key}); err == nil {

					r.hunter = hunter
					r.hunter.SetSize(r.options.Count)

					count := hunter.GetPoints()

					if count == 0 {
						gologger.Error().Msgf("大牛，您的 %s 积分用完了，明天再试试", hunter.DesensitizationKey(key))
						continue
					}

					gologger.Info().Msgf("正在使用 Key: %s, 剩余积分: %d", hunter.DesensitizationKey(key), count)

					return nil
				}
			}

		}
	} else {
		if !r.config.IsFofa() {
			return fmt.Errorf("missing fofa email and key")
		}

		if r.config.IsFofa() {
			fofaEmail, fofaKey := r.config.Fofa.Email, r.config.Fofa.Key
			if fofa, err := fofa.New(&fofa.FofaOptions{Email: fofaEmail, Key: fofaKey}); err == nil {
				r.fofa = fofa
				r.fofa.SetSize(r.options.Count)
				return nil
			}
		}
	}

	return err
}
