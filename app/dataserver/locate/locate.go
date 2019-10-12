package locate

import (
	"github.com/nats-io/nats.go"
	"goss/pkg/types"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var objs = make(map[string]int)
var mutex sync.Mutex

func IsExist(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func Locate(hash string) int {
	log.Println(objs)
	mutex.Lock()
	id, ok := objs[hash]
	mutex.Unlock()
	if !ok {
		return -1
	}
	return id
}

func Add(hash string, id int) {
	mutex.Lock()
	objs[hash] = id
	mutex.Unlock()
}

func Del(hash string) {
	mutex.Lock()
	delete(objs, hash)
	mutex.Unlock()
}

//check obj exist
func StartLocate() {
	//q := rabbitmq.New(os.Getenv("MQ_SERVER"))
	//defer q.Close()
	//q.Bind("dataServers")
	//c := q.Consume()
	//for msg := range c {
	//	hash, e := strconv.Unquote(string(msg.Body))
	//	log.Println("apiServer check file exist hash:",hash)
	//	if e != nil {
	//		panic(e)
	//	}
	//	id := Locate(hash)
	//	if id!=-1 {
	//		q.Send(msg.ReplyTo, types.LocateMessage{Addr:os.Getenv("LISTEN_ADDR"),Id:id})
	//	}
	//}
	nc, e := nats.Connect(os.Getenv("MQ_SERVER"))
	if e != nil {
		panic(e)
	}
	defer nc.Close()
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	c.Subscribe("objwhere", func(subj, reply string, objname string) {
		//hash, e := strconv.Unquote(objname)
		hash:=objname

		log.Println("apiServer check file exist hash:", hash)
		if e != nil {
			log.Println(e)
			return
		}
		id := Locate(hash)
		if id != -1 {
			c.Publish(reply, &types.LocateMessage{Addr: os.Getenv("LISTEN_ADDR"), Id: id})
		}else{
			log.Println("obj is not exist")
		}
	})

	select {}
}

func CollectObjects() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_PATH") + "/objects/*")
	for i := range files {
		file := strings.Split(filepath.Base(files[i]), ".")
		if len(file) != 3 {
			panic(files[i])
		}
		hash := file[0]
		id, e := strconv.Atoi(file[1])
		if e != nil {
			panic(e)
		}
		objs[hash] = id
	}
}
