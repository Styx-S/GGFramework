package Define

type Response struct {
	Ret int
	Msg string
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
	Type     WSPacketType      `json:"type"`
	Category WSPacketCategory  `json:"category"`
	Code     int               `json:"code"`
	Param    map[string]string `json:"params"`
}
