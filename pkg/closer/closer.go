package closer

import (
	"context"
	"example/pkg/logger"
	"os"
	"os/signal"
	"sync"

	"github.com/pkg/errors"
)

var gCloser = New()

func Add(f ...func() error) {
	gCloser.Add(f...)
}

func Wait() {
	gCloser.Wait()
}

func CloseAll() {
	gCloser.CloseAll()
}

type Closer struct {
	m     sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func() error
}

func New(s ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(s) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, s...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}
	return c
}

func (c *Closer) Add(f ...func() error) {
	c.m.Lock()
	c.funcs = append(c.funcs, f...)
	c.m.Unlock()
}

func (c *Closer) Wait() {
	<-c.done
}

func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.m.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.m.Unlock()

		errs := make(chan error, len(funcs))
		for i := range funcs {
			go func(f func() error) {
				errs <- f()
			}(funcs[i])
		}

		for i := 0; i < cap(errs); i++ {
			err := <-errs
			if errClose := errors.Wrap(err, "error from func"); errClose != nil {
				logger.Error(context.Background(), errClose.Error)
			}

		}
	})
}
