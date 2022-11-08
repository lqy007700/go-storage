package locate

import (
	"go-storage/lib/rabbitmq"
	"os"
	"strconv"
	"time"
)

func Locate(name string) string {
	mq := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))

	mq.Publish("dataServers", name)
	c := mq.Consume()

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
