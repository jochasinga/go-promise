package promise

import "fmt"

// Thenable has a Then method
type Thenable interface {
	Then(ResolveFunc, ...RejectFunc) Thenable
}

// Promise represents resolved and rejected channels
type Promise struct {
	resolved chan interface{}
	rejected chan error
}

type (
	// ResolveFunc represents a func(interface{})
	ResolveFunc func(interface{})

	// RejectFunc represents a func(error)
	RejectFunc func(error)
)

// New constructs a Promise around a function which returns an interface{} and error type
func New(fn func() (interface{}, error)) *Promise {
	p := &Promise{
		resolved: make(chan interface{}, 1),
		rejected: make(chan error, 1),
	}
	go func() {
		if res, err := fn(); err != nil {
			p.rejected <- err
		} else {
			p.resolved <- res
		}
	}()
	return p
}

func NewPromise() *Promise {
	return &Promise{
		resolved: make(chan interface{}, 1),
		rejected: make(chan error, 1),
	}
}

// From creates a new Promise from a resolve channel and optional error channel
func From(rc chan interface{}, errc ...chan error) *Promise {
	p := &Promise{}
	return p.From(rc, errc...)
}

// From makes an existing Promise work from a resolve and optional error channel
func (p *Promise) From(rc chan interface{}, errc ...chan error) *Promise {
	if rc != nil {
		p.resolved = rc
	}
	if len(errc) > 0 {
		p.rejected = errc[0]
	}
	return p
}

// To convert a Promise to resolve and reject channels
func (p *Promise) To() (chan interface{}, chan error) {
	rc := make(chan interface{})
	errc := make(chan error)
	go func() {
		if result, ok := <-p.resolved; ok {
			rc <- result
			close(rc)
		}
	}()
	go func() {
		if err, ok := <-p.rejected; ok {
			errc <- err
			close(errc)
		}
	}()
	return rc, errc
}

// Then accepts a ResolveFunc and an optional RejectFunc to handle future result
func (p *Promise) Then(resolve ResolveFunc, reject ...RejectFunc) Thenable {
	go func() {
		select {
		case result := <-p.resolved:
			if resolve != nil {
				resolve(result)
			}
		case err := <-p.rejected:
			if len(reject) > 0 {
				if reject[0] != nil {
					reject[0](err)
				}
			}
		}
	}()
	return p
}

// Resolve saves the interface{} value for next Then chain
func (p *Promise) Resolve(any interface{}) Thenable {
	if _, ok := <-p.resolved; !ok {
		p.resolved = make(chan interface{}, 1)
	}
	go func() {
		p.resolved <- any
	}()
	return p
}

// Reject saves the error value for next Then chain
func (p *Promise) Reject(reason string, values ...interface{}) Thenable {
	if _, ok := <-p.rejected; !ok {
		p.rejected = make(chan error, 1)
	}
	go func() {
		p.rejected <- fmt.Errorf(reason, values...)
	}()
	return p
}
