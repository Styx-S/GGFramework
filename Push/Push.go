package Push

import (
	"GGFramework/Define"
)

type ModuleImpl struct {
	svr   *Define.Server
	Items map[Define.RoomID]*[]*WSContext
	/* 用于在单个线程中执行（可能）是多线程发来的请求 */
	RegisterChan   chan *WSContext
	UnRegisterChan chan *WSContext
	KickChan       chan *WSContext
}

func (impl *ModuleImpl) Broadcast(packet *Define.BroadcastPacket) {
	impl.sendAll(&packet.Packet, packet.RoomID, packet.IgnoreUser)
}

func (impl *ModuleImpl) RemoveNotify(roomID Define.RoomID, username Define.UserName) {

}

func (impl *ModuleImpl) New() *ModuleImpl {
	impl.Items = make(map[Define.RoomID]*[]*WSContext)
	impl.RegisterChan = make(chan *WSContext)
	impl.UnRegisterChan = make(chan *WSContext)
	impl.KickChan = make(chan *WSContext)
	return impl
}
func (impl *ModuleImpl) ConnectToSvr(svr *Define.Server) {
	impl.svr = svr
}

func (impl *ModuleImpl) Start() {
	for {
		select {
		case ctx := <-impl.RegisterChan:
			if impl.addItem(ctx) {

			}

		case ctx := <-impl.UnRegisterChan:
			if impl.removeItem(ctx) {

			}
		}
	}
}

func (impl *ModuleImpl) Notify(roomID Define.RoomID, username Define.UserName, packet *Define.WSPacket) {
	list, ok := impl.Items[roomID]
	if ok {
		for _, item := range *list {
			if item.UserID == username {
				item.Send <- packet
				return
			}
		}
	}
}

func (impl *ModuleImpl) sendAll(packet *Define.WSPacket, roomID Define.RoomID, ignoreUser Define.UserName) bool {
	list, ok := impl.Items[roomID]
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

func (impl *ModuleImpl) addItem(ctx *WSContext) bool {
	list, ok := impl.Items[ctx.RoomID]
	if !ok {
		return false
	}
	*list = append(*list, ctx)
	return true
}

func (impl *ModuleImpl) removeItem(ctx *WSContext) bool {

	_, target, list, ok := impl.findContext(ctx.RoomID, ctx.UserID)

	if ok {
		*list = append((*list)[:target], (*list)[target:]...)
		return true
	}

	return false
}

// target Context, its index, room list,
func (impl *ModuleImpl) findContext(roomID Define.RoomID, userID Define.UserName) (*WSContext, int, *[]*WSContext, bool) {
	list, ok := impl.Items[roomID]
	if ok {
		for index, item := range *list {
			if item.UserID == userID && item.RoomID == roomID {
				return item, index, list, true
			}
		}
	}
	return nil, -1, list, false
}
