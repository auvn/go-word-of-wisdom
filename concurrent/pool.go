package concurrent

import (
	"sync"
)

var _ Goer = (*WorkersPool)(nil)

type WorkersPool struct {
	jobs chan func()
	wg   sync.WaitGroup
}

func NewWorkersPool(size int) *WorkersPool {
	p := WorkersPool{
		jobs: make(chan func()),
	}

	p.wg.Add(size)
	for i := 0; i < size; i++ {
		go func() {
			p.workerLoop()
			p.wg.Done()
		}()
	}

	return &p
}

func (p *WorkersPool) Go(fn func()) {
	p.jobs <- fn
}

func (p *WorkersPool) workerLoop() {
	for fn := range p.jobs {
		fn()
	}
}

func (p *WorkersPool) Close() {
	close(p.jobs)
}

func (p *WorkersPool) Wait() {
	p.wg.Wait()
}
