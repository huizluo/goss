package locate

import (
	"github.com/nats-io/nats.go"
	"goss/pkg/rs"
	"goss/pkg/types"
	"log"
	"os"
	"time"
)

func Locate(name string) (locateInfo map[int]string) {

	nc, e := nats.Connect(os.Getenv("MQ_SERVER"))
	if e != nil {
		return nil
	}
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	locateInfo = make(map[int]string)
	for i := 0; i < rs.ALL_SHARDS; i++ {
		var info types.LocateMessage
		log.Println("check obj where!!!",name)
		e = c.Request("objwhere", name, &info, time.Second*2)
		log.Println(info.Addr)
		if e != nil {
			log.Println("response error ", e)
			return
		}
		locateInfo[info.Id] = info.Addr
	}
	return
}

func IsExist(name string) bool {
	return len(Locate(name)) >= rs.DATA_SHARDS
}
