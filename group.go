package epollevts

import (
	"epollevts/topics"
	"errors"
	"log"

	"golang.org/x/sys/unix"
)

type TopicGroupCallFn func(evt_value []topics.EvtMsg)

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
	msgs := make([]topics.EvtMsg, 0, eb.perProcessTopicCap)
	var buf [8]byte
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
			_, err := unix.Read(int(topic.GetFd()), buf[:])
			//topic.GetLock().Unlock()
			if err != nil {
				log.Printf("Read from eventfd error: %v", err)
			}

			msgs = append(msgs, topics.EvtMsg{
				Key:   topic.GetKey(),
				Count: 1,
				Data:  buf,
			})

			eb.callBack(msgs)
			msgs = msgs[:0]

			//清零复用
			for k, _ := range buf {
				buf[k] = 0
			}
		}
	}
}
