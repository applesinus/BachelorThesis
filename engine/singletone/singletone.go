package singletone

import (
	"BachelorThesis/engine/collision"
	"BachelorThesis/engine/objects"
	"context"
	"sync"
)

type Engine struct {
	Algorithm          string
	SecondaryAlgorithm string
	ResolveAlgorithm   string

	Context        context.Context
	CollisionStart chan struct{}
	CollisionEnd   chan struct{}

	mu *sync.Mutex

	ObjectPool *[]objects.Object
}

func NewEngine(algorithm, secondaryAlgorithm, resolveAlgorithm string, pool *[]objects.Object, ctx context.Context) *Engine {
	return &Engine{
		Algorithm:          algorithm,
		SecondaryAlgorithm: secondaryAlgorithm,
		ResolveAlgorithm:   resolveAlgorithm,

		Context:        ctx,
		CollisionStart: make(chan struct{}),
		CollisionEnd:   make(chan struct{}),

		mu: new(sync.Mutex),

		ObjectPool: pool,
	}
}

func (e *Engine) StartEngineLoop() {
	for {
		select {
		case <-e.Context.Done():
			e.mu.Lock()
			for i := range *e.ObjectPool {
				(*e.ObjectPool)[i] = nil
			}
			e.ObjectPool = nil
			e.mu.Unlock()
			return

		case <-e.CollisionStart:
			e.Mute()
			e.update()
			e.Unmute()
			e.CollisionEnd <- struct{}{}
		}
	}
}

func (e *Engine) update() {
	collision.ProcessCollisions(e.ObjectPool, e.Algorithm, e.SecondaryAlgorithm, e.ResolveAlgorithm)
}

func (e *Engine) AddObject(object objects.Object) {
	e.mu.Lock()
	*e.ObjectPool = append(*e.ObjectPool, object)
	e.mu.Unlock()
}

func (e *Engine) Mute() {
	e.mu.Lock()
}

func (e *Engine) Unmute() {
	e.mu.Unlock()
}

func (e *Engine) ProcessCollisions() {
	e.CollisionStart <- struct{}{}

	<-e.CollisionEnd
}
