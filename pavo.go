package pavo

import (
	"strings"

	"github.com/zan8in/pavo/pkg/config"
	"github.com/zan8in/pavo/pkg/runner"
	sliceutil "github.com/zan8in/pins/slice"
	urlutil "github.com/zan8in/pins/url"
)

var cfg *config.Config

func init() {
	cfg, _ = config.NewConfig()
}

func Query(q []string, size int) ([]string, error) {
	var result []string

	if size == 0 {
		size = 100
	}

	options, err := runner.NewOptions(runner.Options{
		Query: q,
		Count: size,
	})
	if err != nil {
		return result, err
	}

	r, err := runner.NewRunner(options)
	if err != nil {
		return result, err
	}

	if err := r.Run(); err != nil {
		return result, err
	}

	rs := r.Result.GetResult()
	for s := range rs {
		result = append(result, s[0])
	}

	return DedupDomain(result), nil
}

func QuerySubDomain(domain string, size int) ([]string, error) {
	var result []string

	if size == 0 {
		size = 100
	}

	options, err := runner.NewOptions(runner.Options{
		Query: []string{`domain="` + domain + `"`},
		Count: size,
	})
	if err != nil {
		return result, err
	}

	r, err := runner.NewRunner(options)
	if err != nil {
		return result, err
	}

	if err := r.Run(); err != nil {
		return result, err
	}

	rs := r.Result.GetResult()
	for s := range rs {
		result = append(result, s[0])
	}

	return DedupDomain(result), nil
}

func DedupDomain(s []string) []string {
	var n []string
	for _, d := range s {
		if r, err := urlutil.DomainName(d); err == nil {
			n = append(n, strings.TrimSpace(strings.TrimRight(r, ".")))
		}
	}
	return sliceutil.Dedupe(n)
}

type (
	Results struct {
		Result []Result
	}
	Result struct {
		IP       string
		Port     string
		Domain   string
		Protocol string
		Server   string
	}
)

func QueryIPPort(ip string, size int) (Results, error) {
	var result Results

	if size == 0 {
		size = 100
	}

	options, err := runner.NewOptions(runner.Options{
		Query: []string{`ip="` + ip + `"`},
		Count: size,
	})
	if err != nil {
		return result, err
	}

	r, err := runner.NewRunner(options)
	if err != nil {
		return result, err
	}

	if err := r.Run(); err != nil {
		return result, err
	}

	rs := r.Result.GetResult()
	for s := range rs {
		result.Result = append(result.Result, Result{
			IP: s[2], Port: s[3], Domain: s[4], Protocol: s[5], Server: s[6]})
	}

	return result, nil
}

func IsFofa() bool {
	return cfg.IsFofa()
}

func IsHunter() bool {
	return cfg.IsHunter()
}
