package Define

type Response struct {
	Ret   int         `json:"ret"`
	Msg   string      `json:"msg"`
	Param interface{} `json:"param"`
}

type WSPacketType int

const (
	WSPacketTypeHeartbeat    = 0
	WSPacketTypeHeartbeatAck = 1
	WSPacketTypeRawMsg       = 2
)

type WSPacketCategory int

const (
	WSPacketCategoryWebsocket = 0
	WSPacketCategoryRoom      = 1
	WSPacketCategoryGameLogic = 2
)

type WSPacket struct {
	Type     WSPacketType     `json:"type"`
	Category WSPacketCategory `json:"category"`
	Code     int              `json:"code"`
	Param    interface{}      `json:"params"`
}

// 所有websocket code定义在这
const (
	WSCRoomUserJoin        = 0 // 玩家加入房间
	WSCRoomUserLeave       = 1 // 玩家离开房间
	WSCRoomUserChangeReady = 2 // 玩家更改准备状态
	WSCRoomStartGame       = 3 // 所有玩家准备完毕，开始游戏
)

/* 接口有意义的结构体定义在这 */

type UserRsp struct {
	Username   string `json:"username"`
	ReadyState bool   `json:"ready_state"`
}
