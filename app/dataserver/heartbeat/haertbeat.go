package heartbeat

import (
	"github.com/nats-io/nats.go"
	"os"
	"time"
)

func StartHeartbeat() {
	nc, e := nats.Connect(os.Getenv("MQ_SERVER"))
	if e != nil {
		panic(e)
	}
	c, _ := nats.NewEncodedConn(nc, nats.DEFAULT_ENCODER)

	for {
		c.Publish("data_servers", os.Getenv("LISTEN_ADDR"))
		time.Sleep(time.Second * 3)
	}
}
