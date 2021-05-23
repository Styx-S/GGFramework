package Push

import (
	"GGFramework/Define"
	"GGFramework/Room"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

type ConnectWSRequest struct {
	RoomID   string `form:"room_id"`
	Username string `form:"username"`
}

// 处理websocket连接过程
func ConnectWebsocket(ctx *gin.Context) {

	var request ConnectWSRequest
	paramErr := ctx.ShouldBind(&request)
	if paramErr != nil {
		ctx.JSON(http.StatusOK, Define.Response{
			Ret: Define.ParametersError,
			Msg: paramErr.Error(),
		})
	}

	// 升级到websocket协议
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {

	}

	ws := &WSContext{
		RoomID: Room.ID(request.RoomID),
		UserID: Room.UserName(request.Username),
		Socket: conn,
		Send:   make(chan *Define.WSPacket),
	}

	go ws.Write()
	go ws.Read()
	CtxManager.Register <- ws
}

func (ctx *WSContext) Read() {
	defer func() {
		CtxManager.UnRegister <- ctx
		ctx.Socket.Close()
	}()

	for {
		_, message, err := ctx.Socket.ReadMessage()
		if err != nil {
			CtxManager.UnRegister <- ctx
			ctx.Socket.Close()
			break
		}
		var packet Define.WSPacket
		jsonErr := json.Unmarshal(message, &packet)
		if jsonErr != nil {
			continue
		}

		ctx.Dispatch(&packet)
	}
}

func (ctx *WSContext) Write() {
	defer func() {
		ctx.Socket.Close()
	}()

	for {
		select {
		case packet, ok := <-ctx.Send:
			if !ok {
				continue
			}

			message, err := json.Marshal(packet)
			if err != nil {
				continue
			}
			ctx.Socket.WriteMessage(int(packet.Type), message)
		}
	}
}

func (ctx *WSContext) Dispatch(packet *Define.WSPacket) {
	switch packet.Category {
	case Define.WSPacketCategoryWebsocket:
		{
			switch packet.Type {
			case Define.WSPacketTypeHeartbeat:
				ctx.Send <- &Define.WSPacket{
					Type:     Define.WSPacketTypeHeartbeatAck,
					Category: Define.WSPacketCategoryWebsocket,
					Param: map[string]string{
						"duration": strconv.Itoa(10),
					},
				}
			default:
				return
			}
		}
	case Define.WSPacketCategoryGameLogic:
		{

		}

	}
}

// server端不需要做heartbeat计时逻辑，交给client来做
//func (ctx *WSContext) HeartBeat() {
//	for {
//		select {
//		case next := <-ctx.NextHeartBeat:
//			if next {
//				// heartbeat
//				var heartbeatPacket = Define.WSPacket{
//
//				}
//				json.Marshal()
//
//				// next heart beat
//				go func() {
//					//TODO: 当得到stop的消息时似乎没机会通知这个协程？由此是否会导致该ctx延迟释放
//					time.Sleep(10)
//					ctx.NextHeartBeat <- true
//				}()
//			} else  {
//				// stop
//				return
//			}
//
//		}
//	}
//}
