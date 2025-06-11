package worker

import (
	"sync"
)

type WorkerWrapper struct {
	Worker   *Worker
	StopChan chan struct{}
}

type WorkerManager struct {
	mu      sync.Mutex
	Workers map[uint]*WorkerWrapper
}

func NewWorkerManager() *WorkerManager {
	return &WorkerManager{
		Workers: make(map[uint]*WorkerWrapper),
	}
}

func (wm *WorkerManager) Add(ProcessID uint, wrapper *WorkerWrapper) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.Workers[ProcessID] = wrapper
}

func (wm *WorkerManager) Remove(ProcessID uint) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	delete(wm.Workers, ProcessID)
}

func (wm *WorkerManager) Stop(ProcessID uint) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	if wrapper, exists := wm.Workers[ProcessID]; exists {
		close(wrapper.StopChan)
		delete(wm.Workers, ProcessID)
	}
}

func (wm *WorkerManager) List() []uint {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	ids := make([]uint, 0, len(wm.Workers))
	for id := range wm.Workers {
		ids = append(ids, id)
	}
	return ids
}
