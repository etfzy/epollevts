package topics

import (
	"errors"
	"sync"

	"golang.org/x/sys/unix"
)

const (
	LT = unix.EPOLLIN
	ET = unix.EPOLLIN | unix.EPOLLET
)

type Topic interface {
	Close() error
	GetFd() int32
	GetKey() string
	GetMode() uint32
	SetStop()
	GetState() bool
	SetStart()
	Publish() error
	GetLock() *sync.Mutex
}

type EpollManger struct {
	topic       Topic
	epoll_event *unix.EpollEvent
}

type Center struct {
	lock    sync.RWMutex
	mTopics map[string]*EpollManger
	mFds    map[int32]string
}

func CreateTopicsCenter() *Center {
	return &Center{
		lock:    sync.RWMutex{},
		mTopics: map[string]*EpollManger{},
		mFds:    map[int32]string{},
	}
}

func (tc *Center) AddTopic(epfd int, topic Topic) error {
	tc.lock.Lock()
	defer tc.lock.Unlock()
	_, ok := tc.mTopics[topic.GetKey()]
	if ok {
		return errors.New("repeat event key...")
	}

	//加入到epoll 事件中
	ep_evt := unix.EpollEvent{Events: topic.GetMode(), Fd: int32(topic.GetFd())}
	if err := unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, int(topic.GetFd()), &ep_evt); err != nil {
		return err
	}

	tc.mTopics[topic.GetKey()] = &EpollManger{
		topic:       topic,
		epoll_event: &ep_evt,
	}
	tc.mFds[topic.GetFd()] = topic.GetKey()
	return nil
}

func (tc *Center) DelTopic(epfd int, eventKey string) error {
	tc.lock.Lock()
	defer tc.lock.Unlock()
	tp, ok := tc.mTopics[eventKey]
	if !ok {
		return nil
	}

	if err := unix.EpollCtl(epfd, unix.EPOLL_CTL_DEL, int(tp.topic.GetFd()), tp.epoll_event); err != nil {
		return err
	}

	delete(tc.mFds, tp.topic.GetFd())
	delete(tc.mTopics, eventKey)

	return nil
}

func (tc *Center) GetTopicByKey(eventKey string) Topic {
	tc.lock.RLock()
	defer tc.lock.RUnlock()
	_, ok := tc.mTopics[eventKey]
	if !ok {
		return nil
	}
	return tc.mTopics[eventKey].topic
}

func (tc *Center) GetTopicByFd(fd int32) Topic {
	tc.lock.RLock()
	defer tc.lock.RUnlock()
	evtkey, ok := tc.mFds[fd]
	if !ok {
		return nil
	}

	return tc.mTopics[evtkey].topic
}
