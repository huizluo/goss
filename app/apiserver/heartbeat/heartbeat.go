package heartbeat

import (
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"sync"
	"time"
)

var dataServers = make(map[string]time.Time)
var mux sync.Mutex

func ListenHeartbeat() {
	nc, e := nats.Connect(os.Getenv("MQ_SERVER"))
	if e != nil {
		panic(e)
	}
	c, _ := nats.NewEncodedConn(nc, nats.DEFAULT_ENCODER)

	c.Subscribe("data_servers", func(msg *nats.Msg) {
		mux.Lock()
		dataServers[string(msg.Data)] = time.Now()
		mux.Unlock()

	})

	select {}
	//dataServers["127.0.0.1:8061"] = time.Now()
	//dataServers["127.0.0.1:8062"] = time.Now()
	//dataServers["127.0.0.1:8063"] = time.Now()
}

func RemoveExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mux.Lock()
		for s, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
			log.Println("dataserver len:", len(dataServers))
		}
		mux.Unlock()
	}
}

func GetDataServers() []string {
	mux.Lock()
	defer mux.Unlock()
	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	return ds
}
