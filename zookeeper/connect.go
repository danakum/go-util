package zookeeper

import (
	"github.com/danakum/go-util/log"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

func Connect() *zk.Conn {
	c, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
	if err != nil {
		log.Fatal(err)
	}

	return c
}
