package handles

import (
	"github.com/gorilla/websocket"
	"github.com/limitedlee/microservice/common/socket"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		// 读取存储空间大小
		ReadBufferSize: 1024,
		// 写入存储空间大小
		WriteBufferSize: 1024,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	WebsocketConns = map[string]*socket.Connection{}
)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	var (
		wsConn *websocket.Conn
		err    error
		// data []byte
		conn *socket.Connection
		data []byte
	)
	// 完成http应答，在httpheader中放下如下参数
	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return // 获取连接失败直接返回
	}
	if conn, err = socket.InitConnection(wsConn); err != nil {
		goto ERR
	}

	//未设置连接账号信息
	if r.URL.RawQuery == "" {
		return
	}
	//每次都使用个最新的连接池
	WebsocketConns[r.URL.RawQuery] = conn

	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		println("当前连接的对象编号：", r.URL.RawQuery)
		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}
	}

ERR:
	// 关闭当前连接
	WebsocketConns[r.URL.RawQuery].Close() //关闭当前连接
}
