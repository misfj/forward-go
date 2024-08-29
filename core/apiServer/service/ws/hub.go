package ws

import (
	"forward-go/global"
	"forward-go/log"
	"github.com/gorilla/websocket"
	"time"
)

type MyWsConn struct {
	conn             *websocket.Conn
	LastTimeKeepLive time.Time
}

func (c *MyWsConn) ReadData() {
	for {
		//检查上次保活时间距离现在时长,你可能使用conn.SetReadDeadline(长时间反复测试,有点问题)或者SetPingHandler进行心跳保活
		if time.Now().Sub(c.LastTimeKeepLive) >
			time.Duration(global.GlobalConfig.Forward.HealthyInterval)*time.Minute {
			//主动断开连接
			err := c.conn.Close()
			if err != nil {
				log.Error(err)
				return
			}
		}
		//一直读取数据
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Error(err)
			return
		}
		wrapMsg := &wrapMsg{
			msgID:      "",
			msg:        nil,
			msgSrc:     "",
			msgDes:     "",
			msgRecount: 0,
			msgErr:     nil,
		}
	}

}
