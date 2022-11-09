package heartbeat

import (
	"go-storage/lib/rabbitmq"
	"os"
	"time"
)

// StartHeartbeat 心跳 获取本机服务地址发送消息队列
func StartHeartbeat() {
	mq := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer mq.Close()

	for {
		mq.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}
