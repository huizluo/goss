package heartbeat

import (
	"goss/app/dataserver/objects"
	"goss/pkg/rabbitmq"
	"os"
)

func StartHeartbeat() {
	q := rabbitmq.New(os.Getenv("MQ_SERVER"))
	defer q.Close()

	for {
		q.Publish("apiServers", objects.LISTEN_ADDRESS)
	}
}
