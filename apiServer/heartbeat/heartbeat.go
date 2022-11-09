package heartbeat

import (
	"go-storage/common"
	"go-storage/lib/rabbitmq"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

// 数据服务节点
var dataServers = make(map[string]time.Time)
var mutex = sync.RWMutex{}

// ListenHeartbeat 监听服务心跳
func ListenHeartbeat() {
	mq := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer mq.Close()
	mq.Bind("apiServers")

	c := mq.Consume()
	go removeExpiredDataServer()

	// 接收心跳包
	for msg := range c {
		addr, err := strconv.Unquote(string(msg.Body))
		common.FailOnError(err, "Unquote fail")

		mutex.Lock()
		dataServers[addr] = time.Now()
		mutex.Unlock()
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