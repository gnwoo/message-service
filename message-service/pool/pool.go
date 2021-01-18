package pool

import (
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)
const (
	// Time allowed to write a message to the peer.
	writeWait = 2 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 5 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
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

func init()  {
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
	newConn := &WsConn{
		id: id,
		c:  c,
	}
	WsConnPool.Conns[id] = newConn
	return newConn
}

func (p *WsPool) RemoveConn(c *WsConn) {
	fmt.Printf("Remove Conn %v\n", c.id)
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
			fmt.Println("Connection corrupted")
			break
		}
		fmt.Println(string(msgVal))
	}
}

func (c *WsConn)Write(msgType int, payload []byte) error {
	c.c.SetWriteDeadline(time.Now().Add(writeWait))
	return c.c.WriteMessage(msgType, payload)
}


//func (p *WsPool) SendTo(id string, msg string) {
//	p.Conns[id].Write(websocket.TextMessage, []byte(msg))
//}