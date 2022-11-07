package heartbeat

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

func StartHeartbeat() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "conn rabbitmq")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "rabbitmq channel")
	defer ch.Close()

	// 3. 声明消息要发送到的队列
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	for {
		err = ch.Publish(
			"apiServers", // exchange
			q.Name,       // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("127.0.0.1:8888"),
			})
		time.Sleep(5 * time.Second)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
