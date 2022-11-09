package locate

import (
	"go-storage/lib/rabbitmq"
	"os"
	"strconv"
	"time"
)

// Locate 发送消息查询是否有对应的文件，有则接收服务地址
func Locate(name string) string {
	mq := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))

	mq.Publish("dataServers", name)
	c := mq.Consume()

	// 超时时间
	go func() {
		time.Sleep(time.Second)
		mq.Close()
	}()

	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

func Exist(name string) bool {
	return Locate(name) != ""
}
