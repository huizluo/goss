package locate

import (
	"goss/pkg/rabbitmq"
	"strconv"
	"time"
)

func Locate(name string) string {
	q := rabbitmq.New("amqp://admin:admin@10.12.32.51:5672")
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

func IsExist(name string) bool {
	return Locate(name) != ""
}
