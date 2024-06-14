package wpool

import (
	"sync"
)

type Pool interface {
	Execute(fn Executive)
	Stop()
	AwaitError() chan error
}

var _ Pool = &pool{}

func NewPool(cap int) *pool {
	p := &pool{
		queueCh: make(chan Executive, cap),
		errorCh: make(chan error, cap),
		stopCh:  make(chan struct{}, 1),
		cap:     cap,
	}

	p.run()

	return p
}

type Executive func() error

type pool struct {
	queueCh chan Executive
	errorCh chan error
	stopCh  chan struct{}

	cap int
}

func (p *pool) run() {
	wg := sync.WaitGroup{}
	wg.Add(p.cap)

	for range p.cap {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			for fn := range p.queueCh {
				if err := fn(); err != nil {
					p.errorCh <- err
				}
			}
		}(&wg)
	}

	go func() {
		wg.Wait()
		close(p.errorCh)
	}()
}

func (p *pool) Execute(fn Executive) {
	select {
	case p.queueCh <- fn:
	}
}

func (p *pool) Stop() {
	close(p.queueCh)
}

func (p *pool) AwaitError() chan error {
	return p.errorCh
}
