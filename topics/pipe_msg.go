package topics

import (
	"errors"
	"fmt"
	"log"
	"os"

	"golang.org/x/sys/unix"
)

type TopicPipe struct {
	Base
	fr  *os.File
	fw  *os.File
	Fwd int32
}

// mode 模式：ET 高性能；LT 不丢消息
func CreateTopicMessage(evtkey string) (Topic, error) {
	// 创建一个eventfd
	pipeR, pipeW, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	et := TopicPipe{}
	et.Base.Stop = false
	et.Base.Fd = int32(pipeR.Fd())
	et.Base.Key = evtkey
	et.Base.Mode = LT
	et.fr = pipeR
	et.fw = pipeW
	et.Fwd = int32(pipeW.Fd())
	fmt.Println(et.Fwd, et.Base.Fd)
	if err := unix.SetNonblock(int(et.Fd), true); err != nil {
		log.Fatalf("Failed to set pipe to non-blocking: %v", err)
	}
	return &et, nil
}

func (et *TopicPipe) Publish() error {
	if et.GetState() {
		return errors.New("not stop topic...")
	}
	buf := []byte{1}

	_, err := unix.Write(int(et.Fwd), buf[:])
	if err != nil {
		return err
	}

	return nil
}

func (et *TopicPipe) Close() error {
	et.SetStop()
	et.fr.Close()
	et.fw.Close()
	return nil
}
