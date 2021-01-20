package internal

import (
	"bytes"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"webssh/flx"
)


var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsSsh(c *gin.Context){
	client, err := flx.NewSshClient("","","")
	if err != nil {
		return
	}
	defer client.Close()
	//startTime := time.Now()
	ssConn, err := NewSshConn(24, 48, client)
	if err != nil{
		return
	}
	defer ssConn.Close()
	// after configure, the WebSocket is ok.
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil{
		return
	}
	defer wsConn.Close()

	quitChan := make(chan bool, 3)

	var logBuff = new(bytes.Buffer)

	// most messages are ssh output, not webSocket input
	go ssConn.ReceiveWsMsg(wsConn, logBuff, quitChan)
	go ssConn.SendComboOutput(wsConn, quitChan)
	go ssConn.SessionWait(quitChan)

	<-quitChan
	logrus.Info("websocket finished")

}