package heartbeat

import (
	"go-storage/lib/rabbitmq"
	"os"
	"time"
)

func StartHeartbeat() {
	mq := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer mq.Close()

	for {
		mq.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}