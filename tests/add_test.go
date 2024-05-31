package epollevts

import (
	"log"
	"testing"
	"time"
)

func TestGroupAdd(t *testing.T) {

	topic_1, err := CreateTopic("test_topic_1", ET)
	if err != nil {
		log.Fatal("create topic failed:", err)
	}

	topic_2, err := CreateTopic("test_topic_2", ET)
	if err != nil {
		log.Fatal("create topic failed:", err)
	}

	//param1:名称
	//param2:控制每次获取最多几个topic的事件
	//param3:回调
	topic_group, err := CreateTopicGroup("topic_group_1", 32, tv)
	if err != nil {
		log.Fatal("create topic group failed:", err)
	}

	//添加topic
	topic_group.AddTopic(topic_1)
	topic_group.AddTopic(topic_2)

	//启动group
	go topic_group.Run()

	//分别对topic1 和topic2 进行消息发送
	go func() {
		for {
			topic_1.Publish()
			time.Sleep(3 * time.Millisecond)
		}

	}()

	go func() {
		for {
			topic_2.Publish()
			time.Sleep(3 * time.Millisecond)
		}
	}()

	time.Sleep(5 * time.Second)

	//测试运行时增加一个topic
	topic_3, err := CreateTopic("test_topic_3", ET)
	if err != nil {
		log.Fatal("create topic failed:", err)
	}

	topic_group.AddTopic(topic_3)
	go func() {
		for {
			topic_3.Publish()
			time.Sleep(3 * time.Millisecond)
		}
	}()

	time.Sleep(5 * time.Second)

}
