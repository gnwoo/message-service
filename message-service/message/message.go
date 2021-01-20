package message

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/go-redis/redis/v8"
	"os"
)

var topic = "C02D3776MD6R"


const ROCKETMQ_NAMESRV = "127.0.0.1:9876"
const ROCKETMQ_GROUP = "testGroup"
const REDIS_ADDR = "localhost:6379"


var MqProducer rocketmq.Producer
var MqConsumer rocketmq.PushConsumer
var RedisClient *redis.Client

func initRocketMQ() {
	var err error
	svr,_:=primitive.NewNamesrvAddr(ROCKETMQ_NAMESRV)
	MqProducer, _ = rocketmq.NewProducer(
		producer.WithNameServer(svr),
		producer.WithRetry(2),
	)
	err = MqProducer.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	MqConsumer,_ = rocketmq.NewPushConsumer(
		consumer.WithGroupName(ROCKETMQ_GROUP),
		consumer.WithNameServer(svr),
	)
	err = MqConsumer.Subscribe(topic, consumer.MessageSelector{}, ConsumeMessageCallback)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		os.Exit(1)
	}
	err = MqConsumer.Start()
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		os.Exit(1)
	}
}

func initRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     REDIS_ADDR,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}


