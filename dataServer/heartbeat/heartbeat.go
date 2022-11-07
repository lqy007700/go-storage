package heartbeat

import (
	"fmt"
	"github.com/Shopify/sarama"
	"go-storage/common"
	"time"
)

func StartHeartbeat() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = "apiServers"
	msg.Value = sarama.StringEncoder("127.0.0.1:8888")
	// 连接kafka
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	common.FailOnError(err, "producer closed, err:")
	defer client.Close()

	for {
		// 间隔5s发送消息
		pid, offset, err := client.SendMessage(msg)
		fmt.Printf("pid:%v offset:%v err:%v \n", pid, offset, err)
		time.Sleep(5 * time.Second)
	}
}