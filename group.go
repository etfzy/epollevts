package epollevts

import (
	"errors"
	"log"
	"syscall"

	"github.com/etfzy/epollevts/topics"

	"golang.org/x/sys/unix"
)

type TopicGroupCallFn func(evt_value []string)

type TopicGroup struct {
	Name               string
	epFd               int
	perProcessTopicCap int //每次获取的事件数量最大值
	topic_center       *topics.Center
	callBack           TopicGroupCallFn
}

func CreateTopicGroup(name string, perProcessTopicCap int, call TopicGroupCallFn) (*TopicGroup, error) {
	epfd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	if perProcessTopicCap == 0 {
		return nil, errors.New("cap must > 0...")
	}
	return &TopicGroup{
		Name:               name,
		epFd:               epfd,
		topic_center:       topics.CreateTopicsCenter(),
		perProcessTopicCap: perProcessTopicCap,
		callBack:           call,
	}, nil
}

func (eb *TopicGroup) Distory() error {
	return unix.Close(eb.epFd)
}

func (eb *TopicGroup) AddTopic(e topics.Topic) error {
	return eb.topic_center.AddTopic(eb.epFd, e)
}

func (eb *TopicGroup) DelTopic(e topics.Topic) error {
	return eb.topic_center.DelTopic(eb.epFd, e.GetKey())
}

func (eb *TopicGroup) GelTopic(key string) topics.Topic {
	return eb.topic_center.GetTopicByKey(key)
}

func (eb *TopicGroup) Run() error {

	events := make([]unix.EpollEvent, eb.perProcessTopicCap)
	msgs := make([]string, 0, eb.perProcessTopicCap)
	var buf [1]byte
	for {
		n, err := unix.EpollWait(eb.epFd, events, -1)
		if err != nil {
			return err
		}

		for i := 0; i < n; i++ {
			evfd := events[i].Fd
			topic := eb.topic_center.GetTopicByFd(evfd)

			//key 不处理
			if topic == nil {
				continue
			}

			//禁用
			if topic.GetState() {
				continue
			}

			//topic.GetLock().Lock()
			frd, _ := topic.GetFd()
			_, err := unix.Read(int(frd), buf[:])
			//topic.GetLock().Unlock()
			if err != nil {
				if err == syscall.EAGAIN {
					continue
				} else {
					log.Printf("Read from eventfd error: %v", err)
				}
			}

			msgs = append(msgs, topic.GetKey())

			eb.callBack(msgs)
			msgs = msgs[:0]

			//清零复用
			buf[0] = 0
		}
	}
}
