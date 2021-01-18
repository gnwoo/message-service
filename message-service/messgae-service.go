package message_service

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/go-redis/redis/v8"
	"gwnoo/message-service/pool"
	"net/http"
	"os"
)


func initRocketMQ() {
	var err error
	svr,_:=primitive.NewNamesrvAddr(ROCKETMQ_NAMESRV)
	p, _ = rocketmq.NewProducer(
		producer.WithNameServer(svr),
		producer.WithRetry(2),
	)
	err = p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	c,_ = rocketmq.NewPushConsumer(
		consumer.WithGroupName(ROCKETMQ_GROUP),
		consumer.WithNameServer(svr),
	)
	err = c.Subscribe(topic, consumer.MessageSelector{}, consumeMessageCallback)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		os.Exit(1)
	}
	err = c.Start()
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		os.Exit(1)
	}
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     REDIS_ADDR,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}


func init() {
	//initRocketMQ()
	//initRedis()
}

func im(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	id := r.Header.Get("whoimi")
	wsC := pool.WsConnPool.NewConn(id,c)

	//bulk, oldOffset, readLength := getOnlineMessages(id)
	//
	//for _,v := range bulk {
	//	wsC.Write(websocket.TextMessage, []byte(v))
	//}
	//offsetKey :=getOffsetKey(id)
	//err = updateConsumeOffset(offsetKey, int64(oldOffset - readLength))
	//if err !=nil {
	//	panic(err)
	//}

	go wsC.ReadConn()
}

func Start() {
	http.HandleFunc("/im", im)
	_ = http.ListenAndServe(*addr, nil)
}