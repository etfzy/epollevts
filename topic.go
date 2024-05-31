package epollevts

import (
	"errors"

	"golang.org/x/sys/unix"
)

type Topic struct {
	Evtfd int32
	Key   string
	Mode  uint32
	Stop  bool
}

const (
	LT = unix.EPOLLIN
	ET = unix.EPOLLIN | unix.EPOLLET
)

// mode 模式：ET 高性能；LT 不丢消息
func CreateTopic(evtkey string, mode uint32) (*Topic, error) {
	// 创建一个eventfd
	evfd, err := unix.Eventfd(0, unix.EFD_NONBLOCK)
	if err != nil {
		return nil, err
	}

	return &Topic{
		Evtfd: int32(evfd),
		Key:   evtkey,
		Mode:  mode,
		Stop:  false,
	}, nil
}

func (et *Topic) Publish() error {
	if et.Stop {
		return errors.New("not stop topic...")
	}
	var buf [8]byte
	buf[0] = 1

	_, err := unix.Write(int(et.Evtfd), buf[:])
	if err != nil {
		return err
	}

	return nil
}

func (et *Topic) Distorys() error {
	et.Stop = true
	return unix.Close(int(et.Evtfd))
}

func (et *Topic) StopUse() {
	et.Stop = true
}

func (et *Topic) RecoverUse() {
	et.Stop = true
}
