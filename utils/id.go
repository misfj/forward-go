package utils

import "sync"

type IDGenerator struct {
	mux    sync.RWMutex
	Number int64
}

func (id *IDGenerator) Now() int64 {
	id.mux.RLock()
	defer id.mux.RUnlock()
	return id.Number
}
func (id *IDGenerator) Gen() int64 {
	id.mux.Lock()
	defer id.mux.Unlock()
	id.Number += 1
	return id.Number
}
func NewIDGenerator(number int64) *IDGenerator {
	return &IDGenerator{Number: number}
}
