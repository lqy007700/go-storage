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

// StartLocate 查看本机是否有对应存储，有则返回本机地址
func StartLocate() {
	mq := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
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
