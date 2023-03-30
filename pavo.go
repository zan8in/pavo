package pavo

import (
	"strings"

	"github.com/zan8in/pavo/pkg/runner"
	sliceutil "github.com/zan8in/pins/slice"
	urlutil "github.com/zan8in/pins/url"
)

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
