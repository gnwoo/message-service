package message_service

import (
	"flag"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

var topic = "C02D3776MD6R"

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var p rocketmq.Producer
var c rocketmq.PushConsumer
var rdb *redis.Client


const ROCKETMQ_NAMESRV = "127.0.0.1:9876"
const ROCKETMQ_GROUP = "testGroup"
const REDIS_ADDR = "localhost:6379"
