package locate

import (
	"github.com/Shopify/sarama"
	"go-storage/common"
	"log"
	"strconv"
	"time"
)

func Locate(name string) string {
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	common.FailOnError(err, "fail to start consumer")

	partitionList, err := consumer.Partitions("dataServers") // 根据topic取到所有的分区
	common.FailOnError(err, "fail to get list of partition")
	log.Println(partitionList)

	// 遍历所有的分区
	for partition := range partitionList {
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("apiServers", int32(partition), sarama.OffsetNewest)
		common.FailOnError(err, "failed to start consumer for partition")

		go func(sarama.PartitionConsumer) {
			time.Sleep(time.Second)
		}(pc)

		msg := <-pc.Messages()

		object, _ := strconv.Unquote(string(msg.Value))
		return object
	}
	return ""
}

func Exist(name string) bool {
	return Locate(name) != ""
}
