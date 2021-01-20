package message_service

import (
	"flag"
	"github.com/gorilla/websocket"
	"gwnoo/logger"
	"gwnoo/message-service/message"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options
var addr = flag.String("addr", "localhost:8080", "http service address")
var sugarlog = logger.Logger.Sugar()

func im(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	id := r.Header.Get("whoimi")
	wsC := message.WsConnPool.NewConn(id,c)

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