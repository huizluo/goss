package heartbeat

import (
	"goss/app/dataserver/objects"
	"goss/pkg/rabbitmq"
)

func StartHeartbeat(){
	q:= rabbitmq.New(objects.RABBITMQ_ADDR)
	defer q.Close()

	for{
		q.Publish("apiServers",objects.LISTEN_ADDRESS)
	}
}