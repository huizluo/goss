package locate

import (
	"goss/app/dataserver/objects"
	"goss/pkg/rabbitmq"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var objs = make(map[string]int)
var mutex sync.Mutex

func IsExist(name string) bool {
	_,err:=os.Stat(name)
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

func StartLocate() {
	q := rabbitmq.New(objects.RABBITMQ_ADDR)
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		hash, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		exist := Locate(hash)
		if exist {
			q.Send(msg.ReplyTo,objects.LISTEN_ADDRESS)
		}
	}
}

func CollectObjects() {
	files, _ := filepath.Glob(objects.STORAGE_PATH + "/objects/*")
	for i := range files {
		hash := filepath.Base(files[i])
		objs[hash] = 1
	}
}