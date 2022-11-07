package locate

import (
	"fmt"
	"github.com/Shopify/sarama"
	"go-storage/common"
	"log"
	"os"
	"strconv"
)

func Locate(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func StartLocate() {
	root := "/Users/lqy007700/Data/storage"

	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	common.FailOnError(err, "fail to start consumer")

	partitionList, err := consumer.Partitions("dataServers") // 根据topic取到所有的分区
	common.FailOnError(err, "fail to get list of partition")

	log.Println(partitionList)
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("dataServers", int32(partition), sarama.OffsetNewest)
		common.FailOnError(err, "failed to start consumer for partition")
		defer pc.AsyncClose()

		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				object, e := strconv.Unquote(string(msg.Value))
				common.FailOnError(e, "异步从每个分区消费信息err")
				if Locate(root + "/objects/" + object) {
					notifySend()
				}
			}
		}(pc)
	}
}

func notifySend() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = "dataServers"
	msg.Value = sarama.StringEncoder("127.0.0.1:8888")
	// 连接kafka
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	common.FailOnError(err, "producer closed, err:")
	defer client.Close()

	// 间隔5s发送消息
	pid, offset, err := client.SendMessage(msg)
	fmt.Printf("pid:%v offset:%v err:%v \n", pid, offset, err)

}
