package topics

import (
	"errors"
	"sync"

	"golang.org/x/sys/unix"
)

// count 是触发的计数
type EvtMsg struct {
	Key   string
	Count uint64
	Data  [8]byte
}

type TopicEvent struct {
	Base
}

// mode 模式：ET 高性能；LT 不丢消息
func CreateTopicEvent(evtkey string, mode uint32) (Topic, error) {
	// 创建一个eventfd
	evfd, err := unix.Eventfd(0, unix.EFD_NONBLOCK)
	if err != nil {
		return nil, err
	}

	et := TopicEvent{}
	et.Base.Stop = false
	et.Base.Fd = int32(evfd)
	et.Base.Key = evtkey
	et.Base.Mode = mode
	et.Base.Llock = &sync.Mutex{}
	return &et, nil
}

func (et *TopicEvent) Publish() error {
	if et.GetState() {
		return errors.New("not stop topic...")
	}
	var buf [8]byte
	buf[0] = 1

	_, err := unix.Write(int(et.GetFd()), buf[:])
	if err != nil {
		return err
	}

	return nil
}

func (et *TopicEvent) Close() error {
	et.SetStop()
	return unix.Close(int(et.GetFd()))
}
