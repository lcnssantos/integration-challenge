package concurrency

import "sync"

type element interface {
	interface{}
}

type PubSub[T element] struct {
	sync.Mutex
	channels map[string]chan T
}

func NewPubSub[T element]() *PubSub[T] {
	return &PubSub[T]{
		channels: make(map[string]chan T),
	}
}

func (p *PubSub[T]) Subscribe(key string) chan T {
	p.Lock()
	defer p.Unlock()
	ch := make(chan T)
	p.channels[key] = ch
	return ch
}

func (p *PubSub[T]) Publish(key string, data T) {
	p.channels[key] <- data
}

func (p *PubSub[T]) Unsubscribe(key string) {
	ch := p.channels[key]
	close(ch)
	delete(p.channels, key)
}
