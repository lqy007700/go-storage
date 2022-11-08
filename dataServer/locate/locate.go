package locate

import (
	"go-storage/common"
	"go-storage/lib/rabbitmq"
	"os"
	"strconv"
)

func Locate(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func StartLocate() {
	mq := rabbitmq.New("127.0.0.1")
	defer mq.Close()
	mq.Bind("dataServers")
	m := mq.Consume()

	for msg := range m {
		object, err := strconv.Unquote(string(msg.Body))
		common.FailOnError(err, "Unquote fail")
		if Locate(os.Getenv("STORAGE_ROOT") + "/objects/" + object) {
			mq.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}