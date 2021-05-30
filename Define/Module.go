package Define

type Module interface {
	// 传递持有该Module的Server实例
	ConnectToSvr(svr *Server)
	// 开始
	Start()
}

const (
	ModuleNamePush = "Push"
	ModuleNameRoom = "Room"
)

/* 推送相关 */
type PushModule interface {
	Module

	// 向指定房间中某个用户发送消息
	Notify(roomID RoomID, username UserName, packet *WSPacket)
	// 广播，参见BroadcastPacket定义
	Broadcast(packet *BroadcastPacket)
	// 后续不再通知某个用户，即断开连接
	RemoveNotify(roomID RoomID, username UserName)
}

/* 房间管理相关 */
type RoomModule interface {
	Module

	// 创建一个房间
	Create(roomID RoomID, firstUser UserName) int
	// 在房间中添加一个玩家
	Join(roomID RoomID, username UserName) int
	// 移除房间中的一个玩家
	Leave(roomID RoomID, user UserName) int
	// 更改房间中一个玩家的准备状态
	ChangeReady(roomID RoomID, user UserName, ready bool) int
	// 获取某个房间的玩家列表
	GetUserList(roomID RoomID) (*[]User, int)
	// 向某房间中所有用户广播一条信息
	BroadCastAll(roomID RoomID, code int, param interface{})
	// 向某房间（除某用户意外）的用户广播一条信息
	BroadCast(roomID RoomID, code int, param interface{}, ignoreName UserName)
}
