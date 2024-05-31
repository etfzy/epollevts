package epollevts

import (
	"errors"
	"sync"

	"golang.org/x/sys/unix"
)

type TopicEpoll struct {
	evt         *Topic
	epoll_event *unix.EpollEvent
}

type TopicsCenter struct {
	lock    sync.RWMutex
	mTopics map[string]*TopicEpoll
	mFds    map[int32]string
}

func createTipicsCenter() *TopicsCenter {
	return &TopicsCenter{
		lock:    sync.RWMutex{},
		mTopics: map[string]*TopicEpoll{},
		mFds:    map[int32]string{},
	}
}

func (tc *TopicsCenter) addTopic(epfd int, e *Topic) error {
	tc.lock.Lock()
	defer tc.lock.Unlock()
	_, ok := tc.mTopics[e.Key]
	if ok {
		return errors.New("repeat event key...")
	}

	//加入到epoll 事件中
	ep_evt := unix.EpollEvent{Events: e.Mode, Fd: int32(e.Evtfd)}
	if err := unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, int(e.Evtfd), &ep_evt); err != nil {
		return err
	}

	tc.mTopics[e.Key] = &TopicEpoll{
		evt:         e,
		epoll_event: &ep_evt,
	}
	tc.mFds[e.Evtfd] = e.Key
	return nil
}

func (tc *TopicsCenter) delTopic(epfd int, eventKey string) error {
	tc.lock.Lock()
	defer tc.lock.Unlock()
	topic, ok := tc.mTopics[eventKey]
	if !ok {
		return nil
	}

	if err := unix.EpollCtl(epfd, unix.EPOLL_CTL_DEL, int(topic.evt.Evtfd), topic.epoll_event); err != nil {
		return err
	}

	delete(tc.mFds, topic.evt.Evtfd)
	delete(tc.mTopics, eventKey)

	return nil
}

func (tc *TopicsCenter) getTopicByKey(eventKey string) *TopicEpoll {
	tc.lock.RLock()
	defer tc.lock.RUnlock()
	_, ok := tc.mTopics[eventKey]
	if !ok {
		return nil
	}
	return tc.mTopics[eventKey]
}

func (tc *TopicsCenter) getTopicByFd(fd int32) *TopicEpoll {
	tc.lock.RLock()
	defer tc.lock.RUnlock()
	evtkey, ok := tc.mFds[fd]
	if !ok {
		return nil
	}

	return tc.mTopics[evtkey]
}
