package locate

import (
	"goss/pkg/rabbitmq"
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

func Locate(hash string) bool {
	mutex.Lock()
	_, ok := objs[hash]
	mutex.Unlock()
	return ok
}

func Add(hash string) {
	mutex.Lock()
	objs[hash] = 1
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
		exist := Locate(hash)
		if exist {
			q.Send(msg.ReplyTo, "127.0.0.1:8061")
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
