package tests

import (
	"errors"
	"github.com/alserov/wpool"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(PoolTestSuite))
}

type PoolTestSuite struct {
	suite.Suite

	pool wpool.Pool
}

func (p *PoolTestSuite) SetupTest() {
	p.pool = wpool.NewPool(4)
}

func (p *PoolTestSuite) TestDefault() {
	err := errors.New("error")

	tasks := []struct {
		fn  wpool.Executive
		err error
	}{
		{
			fn: func() error {
				return nil
			},
			err: nil,
		},
		{
			fn: func() error {
				return nil
			},
			err: nil,
		},
		{
			fn: func() error {
				return err
			},
			err: err,
		},
		{
			fn: func() error {
				return err
			},
			err: err,
		},
		{
			fn: func() error {
				return nil
			},
			err: nil,
		},
		{
			fn: func() error {
				return nil
			},
			err: nil,
		},
		// close
		{
			fn: func() error {
				p.pool.Stop()
				return nil
			},
			err: nil,
		},
	}

	errorsCount := 0

	for _, t := range tasks {
		p.pool.Execute(t.fn)
		if t.err != nil {
			errorsCount++
		}
	}

	for err = range p.pool.AwaitError() {
		p.Require().Error(err)
		errorsCount--
	}

	p.Require().Equal(errorsCount, 0)
}

func (p *PoolTestSuite) TestPanicOnSendOnClosed() {
	defer func() {
		p.Require().NotNil(recover())
	}()

	wp := wpool.NewPool(1)
	wp.Stop()

	wp.Execute(func() error {
		return nil
	})
}
