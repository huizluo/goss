package locate

import (
	"goss/pkg/rabbitmq"
	"log"
	"os"
	"strconv"
	"time"
)

func Locate(name string) string {
	q := rabbitmq.New(os.Getenv("MQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	log.Println("dataserver msg:",s)
	return s
}

func IsExist(name string) bool {
	return Locate(name) != ""
}
