package Push

import (
	"GGFramework/Define"
	"GGFramework/Room"
	"github.com/gorilla/websocket"
)

type WSContext struct {
	RoomID Room.ID
	UserID Room.UserName
	Socket *websocket.Conn
	// channels
	Send chan *Define.WSPacket
}

type BroadcastPacket struct {
	RoomID Room.ID
	Packet Define.WSPacket
}

type Center struct {
	Items      map[Room.ID]*[]*WSContext
	Broadcast  chan *BroadcastPacket
	Register   chan *WSContext
	UnRegister chan *WSContext
}

var CtxManager = Center{
	Items:      make(map[Room.ID]*[]*WSContext),
	Broadcast:  make(chan *BroadcastPacket),
	Register:   make(chan *WSContext),
	UnRegister: make(chan *WSContext),
}

const (
	RoomEventUserJoin = 0
	RoomEventUserLeve = 1
)

func (center *Center) Start() {
	for {
		select {
		case ctx := <-center.Register:
			if center.addItem(ctx) {
				center.SendAll(&Define.WSPacket{
					Type:     Define.WSPacketTypeRawMsg,
					Category: Define.WSPacketCategoryRoom,
					Code:     RoomEventUserJoin,
					Param: map[string]string{
						"username": string(ctx.UserID),
					},
				}, ctx.RoomID, ctx.UserID)
			}

		case ctx := <-center.UnRegister:
			if center.removeItem(ctx) {
				center.SendAll(&Define.WSPacket{
					Type:     Define.WSPacketTypeRawMsg,
					Category: Define.WSPacketCategoryRoom,
					Code:     RoomEventUserLeve,
					Param: map[string]string{
						"username": string(ctx.UserID),
					},
				}, ctx.RoomID, ctx.UserID)
			}

		case broadcast := <-center.Broadcast:
			center.SendAll(&broadcast.Packet, broadcast.RoomID, "")
		}
	}
}

func (center *Center) SendAll(packet *Define.WSPacket, roomID Room.ID, ignoreUser Room.UserName) bool {
	list, ok := center.Items[roomID]
	if !ok {
		return false
	}
	for _, item := range *list {
		if item.UserID != ignoreUser {
			item.Send <- packet
		}
	}
	return true
}

func (center *Center) addItem(ctx *WSContext) bool {
	list, ok := center.Items[ctx.RoomID]
	if !ok {
		return false
	}
	*list = append(*list, ctx)
	return true
}

func (center *Center) removeItem(ctx *WSContext) bool {

	_, target, list, ok := center.findContext(ctx.RoomID, ctx.UserID)

	if ok {
		*list = append((*list)[:target], (*list)[target:]...)
		return true
	}

	return false
}

// target Context, its index, room list,
func (center *Center) findContext(roomID Room.ID, userID Room.UserName) (*WSContext, int, *[]*WSContext, bool) {
	list, ok := center.Items[roomID]
	if ok {
		for index, item := range *list {
			if item.UserID == userID && item.RoomID == roomID {
				return item, index, list, true
			}
		}
	}
	return nil, -1, list, false
}
