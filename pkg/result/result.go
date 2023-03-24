package result

import "sync"

type (
	Result struct {
		sync.RWMutex
		result [][]string
	}
)

func NewResult() *Result {
	return &Result{
		result: make([][]string, 0),
	}
}

func (r *Result) AddResult(rst [][]string) {
	r.Lock()
	defer r.Unlock()

	r.result = append(r.result, rst...)
}

func (r *Result) HasResult() bool {
	r.RLock()
	defer r.RUnlock()

	return len(r.result) > 0
}

func (r *Result) GetResult() chan string {
	r.Lock()

	out := make(chan string)

	go func() {
		defer close(out)
		defer r.Unlock()

		for _, r := range r.result {
			for _, s := range r {
				out <- s
			}
		}
	}()

	return out
}
