package adt

import "sync"

type Barrier struct {
	isClosed bool
	lock *sync.Mutex
	waiter *sync.Cond
}

func InitBarrier(startClosed bool) Barrier {
	b := Barrier{
		isClosed: startClosed,
		lock: &sync.Mutex{},
		waiter: nil,
	}
	b.waiter = sync.NewCond(b.lock)
	return b
}

func (b *Barrier) Wait() {
	b.lock.Lock()
	defer b.lock.Unlock()
	
	for b.isClosed {
		b.waiter.Wait()
	}
}

func (b *Barrier) Open() {
	b.lock.Lock()
	defer b.lock.Unlock()
	
	b.isClosed = false
	b.waiter.Broadcast()
}

func (b *Barrier) Close() {
	b.lock.Lock()
	defer b.lock.Unlock()
	
	b.isClosed = true
}