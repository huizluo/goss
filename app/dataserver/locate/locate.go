package locate

import (
	"goss/app/dataserver/objects"
	"goss/pkg/rabbitmq"
	"os"
	"strconv"
)

func IsExist(name string) bool {
	_,err:=os.Stat(name)
	return !os.IsNotExist(err)
}

func StartLocate()  {
	q:=rabbitmq.New(objects.RABBITMQ_ADDR)
	defer q.Close()
	q.Bind("dataServers")
	c:=q.Consume()
	for msg:=range c{
		obj,e:=strconv.Unquote(string(msg.Body))
		if e!=nil{
			panic(e)
		}
		if IsExist(objects.STORAGE_PATH+"/objects/" + obj){
			q.Send(msg.ReplyTo,objects.LISTEN_ADDRESS)
		}
	}
}