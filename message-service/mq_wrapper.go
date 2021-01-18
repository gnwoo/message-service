package message_service

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"gwnoo/message-service/models"
	"gwnoo/message-service/pool"
)

func consumeMessageCallback(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	for i := range msgs {
		//fmt.Printf("subscribe callback: %v \n", msgs[i])
		parsedMsg := models.Message{}
		jsoniter.Unmarshal(msgs[i].Body, &parsedMsg)

		err := rdb.LPush(context.Background(), parsedMsg.Session, msgs[i].Body).Err()
		if err != nil {
			panic(err)
		}

		if v,ok:=pool.WsConnPool.Conns[parsedMsg.Session];ok{
			v.Write(websocket.TextMessage, msgs[i].Body)
		}else {
			fmt.Println("User offline")
		}
	}

	return consumer.ConsumeSuccess, nil
}

func sendMessage(topic string, body []byte) *primitive.SendResult{
	msg := &primitive.Message{
		Topic: topic,
		Body:  body,
	}
	res, err := p.SendSync(context.Background(), msg)
	if err != nil {
		panic(err)
	}
	return res
}