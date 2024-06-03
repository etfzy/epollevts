package epollevts

import (
	"epollevts/topics"
	"fmt"
	"log"
	"testing"
	"time"
)

func tv1(evt_value []string) {
	for k, _ := range evt_value {
		fmt.Println("topic11111:", evt_value[k])

	}

}

func tv2(evt_value []string) {
	for k, _ := range evt_value {
		fmt.Println("topic2222:", evt_value[k])
	}

}

func tv3(evt_value []string) {
	for k, _ := range evt_value {
		fmt.Println("topic3333:", evt_value[k])
	}

}

func TestGroupParallel(t *testing.T) {
	topic_1, err := topics.CreateTopicMessage("/testtopic_1")
	if err != nil {
		log.Fatal("create topic failed:", err)
	}

	topic_2, err := topics.CreateTopicMessage("/canying")
	if err != nil {
		log.Fatal("create topic failed:", err)
	}

	//启动两个group监听同一个fd
	topic_group1, err := CreateTopicGroup("/topic_group_1", 32, tv1)
	if err != nil {
		log.Fatal("create topic group failed:", err)
	}

	//添加topic
	topic_group1.AddTopic(topic_1)

	//启动group
	go topic_group1.Run()

	topic_group2, err := CreateTopicGroup("topic_group_1", 32, tv2)
	if err != nil {
		log.Fatal("create topic group failed:", err)
	}

	//添加topic
	topic_group2.AddTopic(topic_1)

	//启动group
	go topic_group2.Run()

	topic_group3, err := CreateTopicGroup("topic_group_1", 32, tv3)
	if err != nil {
		log.Fatal("create topic group failed:", err)
	}

	//添加topic
	topic_group3.AddTopic(topic_1)
	topic_group3.AddTopic(topic_2)

	//启动group
	go topic_group3.Run()

	//启动两个线程进行同时写
	sends := 100000
	go func() {
		for i := 0; i < sends; i++ {
			topic_1.Publish()
		}
	}()

	go func() {
		for i := 0; i < sends; i++ {
			topic_1.Publish()
		}
	}()

	go func() {
		for i := 0; i < sends; i++ {
			topic_2.Publish()
		}
	}()

	time.Sleep(5 * time.Second)

}
