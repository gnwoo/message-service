package message

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	jsoniter "github.com/json-iterator/go"
	"gwnoo/message-service/models"
)

func ConsumeMessageCallback(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

	for i := range msgs {
		//fmt.Printf("subscribe callback: %v \n", msgs[i])
		parsedMsg := models.Message{}
		jsoniter.Unmarshal(msgs[i].Body, &parsedMsg)

		//err := RedisClient.LPush(context.Background(), parsedMsg.Session, msgs[i].Body).Err()
		//if err != nil {
		//	panic(err)
		//}
		err:= WsConnPool.SendTo(parsedMsg.Session,  msgs[i].Body)
		if err != nil{
			sugarlog.Infow(err.Error(), "user id", parsedMsg.Session)
		}
	}

	return consumer.ConsumeSuccess, nil
}

func SendChatMessage(body []byte) *primitive.SendResult{
	msg := &primitive.Message{
		Topic: topic,
		Body:  body,
	}
	res, err := MqProducer.SendSync(context.Background(), msg)
	if err != nil {
		panic(err)
	}
	return res
}