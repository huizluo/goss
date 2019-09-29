package heartbeat

import (
	"sync"
	"time"
)

var dataServers = make(map[string]time.Time)
var mux sync.Mutex

func ListenHeartbeat() {
	//q := rabbitmq.New(os.Getenv("MQ_SERVER"))
	//defer q.Close()
	//q.Bind("apiServers")
	//c := q.Consume()
	//go removeExpiredDataServer()
	//for msg := range c {
	//	dataServer, e := strconv.Unquote(string(msg.Body))
	//	if e != nil {
	//		panic(e)
	//	}
	//
	//	mux.Lock()
	//	dataServers[dataServer] = time.Now()
	//	mux.Unlock()
	//}
	dataServers["127.0.0.1:8061"] = time.Now()
	dataServers["127.0.0.1:8062"] = time.Now()
	dataServers["127.0.0.1:8063"] = time.Now()
}

func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mux.Lock()
		for s, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
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
