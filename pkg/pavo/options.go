package pavo

import (
	"errors"
	"strings"

	"github.com/zan8in/goflags"
	"github.com/zan8in/pavo/pkg/result"
)

type OnResultCallback func(result.Result)

type Options struct {
	Query    goflags.StringSlice
	Platform string

	Page  int
	Size  int
	Count int

	OnResult OnResultCallback
}

func NewOptions(options Options) (*Options, error) {
	if err := options.validateOptions(); err != nil {
		return nil, err
	}
	return &options, nil
}

func ParseOptions() (*Options, error) {
	options := &Options{}

	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`Pavo`)

	flagSet.CreateGroup("input", "Input",
		flagSet.StringSliceVarP(&options.Query, "query", "q", nil, "query conditions (comma-separated)", goflags.NormalizedStringSliceOptions),
		flagSet.StringVarP(&options.Platform, "platform", "p", "fofa", "cyberspace mapping platform, support format: fofa,hunter"),
	)

	flagSet.CreateGroup("optimization", "Optimization",
		flagSet.IntVar(&options.Count, "count", DefaultQueryCount, "query count"),
	)

	_ = flagSet.Parse()

	if err := options.validateOptions(); err != nil {
		return nil, err
	}

	return options, nil
}

func (options *Options) validateOptions() (err error) {

	if options.Query == nil {
		return errors.New("no query provided")
	}

	if len(options.Platform) == 0 {
		options.Platform = DefaultPlatform
	} else if len(options.Platform) > 0 && options.Platform != FofaPlatform && options.Platform != HunterPlatform {
		options.Platform = DefaultPlatform
	} else {
		options.Platform = strings.ToLower(options.Platform)
	}

	return err
}
