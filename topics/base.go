package topics

import "sync"

type Base struct {
	Fwd   int32
	Frd   int32
	Key   string
	Mode  uint32
	Stop  bool
	Llock *sync.Mutex
}

func (et *Base) SetStop() {
	et.Stop = true
}

func (et *Base) SetStart() {
	et.Stop = false
}

func (et *Base) GetState() bool {
	return et.Stop
}

func (et *Base) GetFd() (int32, int32) {
	return et.Frd, et.Fwd
}

func (et *Base) GetKey() string {
	return et.Key
}

func (et *Base) GetMode() uint32 {
	return et.Mode
}

func (et *Base) GetLock() *sync.Mutex {
	return et.Llock
}
