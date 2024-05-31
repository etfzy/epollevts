package epollevts

import (
	"encoding/binary"
	"errors"
	"log"

	"golang.org/x/sys/unix"
)

type TopicGroupCallFn func(evt_value []EvtMsg)

type TopicGroup struct {
	Name               string
	epFd               int
	perProcessTopicCap int //每次获取的事件数量最大值
	topic_center       *TopicsCenter
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
		topic_center:       createTipicsCenter(),
		perProcessTopicCap: perProcessTopicCap,
		callBack:           call,
	}, nil
}

func (eb *TopicGroup) Distory() error {
	return unix.Close(eb.epFd)
}

func (eb *TopicGroup) AddTopic(e *Topic) error {
	return eb.topic_center.addTopic(eb.epFd, e)
}

func (eb *TopicGroup) DelTopic(e *Topic) error {
	return eb.topic_center.delTopic(eb.epFd, e.Key)
}

func (eb *TopicGroup) GelTopic(topic string) *Topic {
	return eb.topic_center.getTopicByKey(topic).evt
}

func (eb *TopicGroup) Run() error {

	events := make([]unix.EpollEvent, eb.perProcessTopicCap)
	msgs := make([]EvtMsg, 0, eb.perProcessTopicCap)
	var buf [8]byte
	for {
		n, err := unix.EpollWait(eb.epFd, events, -1)
		if err != nil {
			return err
		}

		for i := 0; i < n; i++ {
			evfd := events[i].Fd
			topic := eb.topic_center.getTopicByFd(evfd)

			//key 不处理
			if topic == nil {
				continue
			}

			//禁用
			if topic.evt.Stop {
				continue
			}

			_, err := unix.Read(int(topic.evt.Evtfd), buf[:])
			if err != nil {
				log.Printf("Read from eventfd error: %v", err)
			}

			msgs = append(msgs, EvtMsg{
				key:   topic.evt.Key,
				count: binary.LittleEndian.Uint64(buf[:]),
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
