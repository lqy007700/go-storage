package heartbeat

import (
	"github.com/Shopify/sarama"
	"go-storage/common"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// 数据服务节点
var dataServers = make(map[string]time.Time)
var mutex = sync.RWMutex{}

// ListenHeartbeat 监听服务心跳
func ListenHeartbeat() {
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	common.FailOnError(err, "fail to start consumer")

	partitionList, err := consumer.Partitions("apiServers") // 根据topic取到所有的分区
	common.FailOnError(err, "fail to get list of partition")
	log.Println(partitionList)

	go removeExpiredDataServer()

	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("apiServers", int32(partition), sarama.OffsetNewest)
		common.FailOnError(err, "failed to start consumer for partition")
		defer pc.AsyncClose()

		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				object, e := strconv.Unquote(string(msg.Value))
				common.FailOnError(e, "异步从每个分区消费信息")

				mutex.Lock()
				dataServers[object] = time.Now()
				mutex.Unlock()
			}
		}(pc)
	}
}

// 移除10s内未收到的心跳节点
func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)

		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}

// GetDataServers 获取服务地址
func GetDataServers() []string {
	mutex.RLock()
	defer mutex.RUnlock()

	ds := make([]string, 0)

	for s := range dataServers {
		ds = append(ds, s)
	}
	return ds
}

// ChooseRandomDataServer 随机选一个服务地址
func ChooseRandomDataServer() string {
	ds := GetDataServers()
	n := len(ds)
	if n == 0 {
		return ""
	}

	return ds[rand.Intn(n)]
}