package errgroup

import (
	"context"
	"strings"
	"sync"
)

// Error is a custom error type to track multiple errors.
type Error struct {
	errs []error
}

// Error returns a string concatenation of multiple errors.
func (e Error) Error() string {
	msgs := make([]string, 0, len(e.errs))
	for _, err := range e.errs {
		if err != nil {
			msgs = append(msgs, err.Error())
		}
	}
	return strings.Join(msgs, ";")
}

// Errors returns the error of the Go method as a slice.
func (e Error) Errors() []error {
	return e.errs
}

// Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero Group is valid and does not cancel.
type Group struct {
	wg     sync.WaitGroup
	mu     sync.Mutex
	cancel func()
	once   sync.Once
	errs   []error
}

// WithContext returns a new Group and an associated Context derived from ctx.
func WithContext(ctx context.Context) (*Group, context.Context) {
	egctx, cancel := context.WithCancel(ctx)
	return &Group{
		cancel: cancel,
	}, egctx
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the Error object if they returns an error.
func (g *Group) Wait() error {
	g.wg.Wait()

	if len(g.errs) > 0 {
		return Error{
			errs: g.errs,
		}
	}

	return nil
}

// Go calls the given function in a new goroutine.
func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.once.Do(func() {
				if g.cancel != nil {
					g.cancel()
				}
			})
			g.mu.Lock()
			if g.errs == nil {
				g.errs = make([]error, 0)
			}
			g.errs = append(g.errs, err)
			g.mu.Unlock()
		}
	}()
}
