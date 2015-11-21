package main

import (
	"sync"
)

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (wgw *WaitGroupWrapper) Wrap(callback func()) {
	wgw.Add(1)
	go func() {
		callback()
		wgw.Done()
	}()
}
