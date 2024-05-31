package epollevts

import (
	"log"
	"testing"
	"time"
)

func TestGroupParallel(t *testing.T) {

	topic_1, err := CreateTopic("test_topic_1", ET)
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

	//启动group
	go topic_group.Run()

	//启动两个线程进行同时写
	go func() {
		for {
			topic_1.Publish()
		}

	}()

	go func() {
		for {
			topic_1.Publish()
		}

	}()

	time.Sleep(5 * time.Second)

}
