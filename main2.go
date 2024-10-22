package main

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// 复用
type pool struct {
	mutex   *sync.Mutex
	cond    *sync.Cond
	workers *stack
	capping int32
	running int32
	waiting int32
}

func newPool(capping int32) *pool {
	mutex := &sync.Mutex{}
	cond := sync.NewCond(mutex)
	ws := &stack{
		cap:  int(capping),
		data: make([]*worker, capping),
	}
	return &pool{mutex: mutex, cond: cond, workers: ws, capping: capping}
}

func (p *pool) addRunning(delta int) int {
	return int(atomic.AddInt32(&p.running, int32(delta)))
}

func (p *pool) addWaiting(delta int) {
	atomic.AddInt32(&p.waiting, int32(delta))
}

type stack struct {
	data  []*worker
	size  int
	index int
	cap   int
}

func (s *stack) push(w *worker) bool {
	if s.size == s.cap {
		return false
	}
	s.data[s.index] = w
	s.index++
	s.size++
	return true
}
func (s *stack) pop() *worker {
	if s.size == 0 {
		return nil
	}
	s.index--
	s.size--
	return s.data[s.index]
}
func (s *stack) isEmpty() bool {
	return s.size == 0

}

type worker struct {
	pool         *pool
	ch           chan func()
	lastUsedTime int64
}

func (w *worker) run() {
	w.pool.addRunning(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Fatal(err)
			}
			w.pool.mutex.Lock()
			defer w.pool.mutex.Unlock()
			w.pool.addRunning(-1)
			w.pool.cond.Signal()
		}()
		for f := range w.ch {
			if f == nil {
				return
			}
			f()
			w.pool.revworkers(w)
		}
	}()
}

func (p *pool) revworkers(w *worker) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if w.pool.workers.push(w) {
		w.pool.cond.Signal()
		return
	}
	w.pool.cond.Broadcast()
	return
}

func (p *pool) retresworkers() *worker {
	p.mutex.Lock()
retry:
	if p.capping == -1 || p.capping > p.running {
		p.mutex.Unlock()
		w := &worker{
			pool: p,
			ch:   make(chan func()),
		}
		p.workers.push(w)
		p.addRunning(1)
		w.run()
		return w
	}
	if w := p.workers.pop(); w != nil {
		p.mutex.Unlock()
		return w
	}

	p.addWaiting(1)
	p.cond.Wait() // block and wait for an available worker
	p.addWaiting(-1)
	goto retry
}

func (p *pool) submit(fn func()) {
	w := p.retresworkers()
	w.ch <- fn
}

func main() {
	st := time.Now().UnixMilli()
	defer func() {
		ut := time.Now().UnixMilli()
		log.Printf("%d ms", ut-st)
	}()
	p := newPool(100)
	wg := sync.WaitGroup{}
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		p.submit(func() {
			defer wg.Done()
			println("hello world")
		})
	}
	wg.Wait()
}
