package wpool

import (
	"context"
	"sync"
)

type Pool interface {
	Execute(fn Executive)
	Stop()
	AwaitError() chan error
}

var _ Pool = &pool{}

func NewPool(cap int64) *pool {
	ctx, cancel := context.WithCancel(context.Background())
	p := &pool{
		queueCh: make(chan Executive, cap),
		errorCh: make(chan error, cap),
		cancel:  cancel,
		ctx:     ctx,
		cap:     cap,
	}

	p.run()

	return p
}

type Executive func() error

type pool struct {
	queueCh chan Executive
	errorCh chan error

	cap int64

	ctx    context.Context
	cancel context.CancelFunc
}

func (p *pool) run() {
	wg := sync.WaitGroup{}
	wg.Add(int(p.cap))

	for range p.cap {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			for {
				select {
				case <-p.ctx.Done():
					return
				case fn := <-p.queueCh:
					if err := fn(); err != nil {
						p.errorCh <- err
					}
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
	case <-p.ctx.Done():
		close(p.queueCh)
	case p.queueCh <- fn:
	}
}

func (p *pool) Stop() {
	p.cancel()
}

func (p *pool) AwaitError() chan error {
	return p.errorCh
}
