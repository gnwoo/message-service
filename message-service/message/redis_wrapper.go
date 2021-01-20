package message

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

func getOffsetKey(deviceId string) string {
	return deviceId+":offset"
}

func getConsumeOffset(key string) int {
	//rdb.SetNX(context.Background(), key, 1, time.Second)
	res,err := RedisClient.Get(context.Background(), key).Result()
	if err != nil {
		panic(err)
	}
	intRes, err := strconv.Atoi(res)
	if err != nil {
		panic(err)
	}
	return intRes
}

func updateConsumeOffset(key string, newOffset int64) error {
	fmt.Printf("Setting %s\n", key)
	return RedisClient.Set(context.Background(), key, newOffset, time.Hour).Err()
}

func getOnlineMessages(deviceId string) ([]string, int, int) {
	offsetKey :=getOffsetKey(deviceId)
	offset :=getConsumeOffset(offsetKey)
	res,err := RedisClient.LRange(context.Background(),deviceId,0,int64(offset)).Result()
	if err != nil {
		panic(err)
	}
	l:=len(res)

	return res, offset, l
}