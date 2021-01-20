package message

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"gwnoo/logger"
	"time"
)

const (
	writeWait = 2 * time.Second
	pongWait = 40 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 512
)

type WsConn struct {
	id string
	c *websocket.Conn
}

type WsPool struct {
	HeartbeatTicker *time.Ticker
	Conns map[string]*WsConn
}
var WsConnPool = WsPool{
		HeartbeatTicker: time.NewTicker(time.Second * 60),
		Conns:             make(map[string]*WsConn),
	}

var heartbeatTicker *time.Ticker
var sugarlog = logger.Logger.Sugar()
func init()  {
	initRocketMQ()
	//initRedis()
	heartbeatTicker = time.NewTicker(pingPeriod)
	go runHeartbeat()
}

func runHeartbeat() {
	for {
		select {
		case t := <-heartbeatTicker.C:
			fmt.Printf("Heartbeat wheel at %v\n", t.String())
			for _,v := range WsConnPool.Conns {
				v.Ping()
			}
		}
	}
}

func (c *WsConn)Ping() {
	c.c.WriteMessage(websocket.PingMessage,[]byte{})
}

func (p *WsPool) NewConn(id string, c *websocket.Conn) *WsConn {
	sugarlog.Infow("New connection", "User id", id)
	newConn := &WsConn{
		id: id,
		c:  c,
	}
	WsConnPool.Conns[id] = newConn
	return newConn
}

func (p *WsPool) RemoveConn(c *WsConn) {
	sugarlog.Infow("Remove Conn", "User id",c.id)
	err := c.c.Close()
	if err != nil {
		panic(err)
	}
	delete(p.Conns, c.id)
}

func (c *WsConn)ReadConn() {
	defer func() {
		WsConnPool.RemoveConn(c)
	}()

	c.c.SetReadLimit(maxMessageSize)
	c.c.SetReadDeadline(time.Now().Add(pongWait))
	c.c.SetPongHandler(func(string) error {
		c.c.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msgVal, err := c.c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				panic(err)
			}
			sugarlog.Infow("Connection corrupted", "user id", c.id)
			break
		}
		sugarlog.Debugw("Received message", "body", string(msgVal))
		SendChatMessage(msgVal)
	}
}

func (c *WsConn)write(msgType int, payload []byte) error {
	c.c.SetWriteDeadline(time.Now().Add(writeWait))
	return c.c.WriteMessage(msgType, payload)
}


func (p *WsPool) SendTo(id string, msg []byte) error {
	if v,ok := p.Conns[id]; ok {
		v.write(websocket.TextMessage, msg)
		return nil
	}
	return errors.New("user not found in connection pool")
}