package locate

import (
	"goss/pkg/rabbitmq"
	"goss/pkg/types"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var objs = make(map[string]int)
var mutex sync.Mutex

func IsExist(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func Locate(hash string) int {
	mutex.Lock()
	id, ok := objs[hash]
	mutex.Unlock()
	if !ok{
		return -1
	}
	return id
}

func Add(hash string,id int) {
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
	q := rabbitmq.New(os.Getenv("MQ_SERVER"))
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		hash, e := strconv.Unquote(string(msg.Body))
		log.Println("apiServer check file exist hash:",hash)
		if e != nil {
			panic(e)
		}
		id := Locate(hash)
		if id!=-1 {
			q.Send(msg.ReplyTo, types.LocateMessage{Addr:os.Getenv("LISTEN_ADDR"),Id:id})
		}
	}
}

func CollectObjects() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_PATH") + "/objects/*")
	for i := range files {
		hash := filepath.Base(files[i])
		objs[hash] = 1
	}
}
