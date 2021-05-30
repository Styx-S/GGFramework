package Room

import (
	"GGFramework/Define"
	"strconv"
)

type ModuleImpl struct {
	svr      *Define.Server
	id2rooms map[Define.RoomID]*Define.Room
}

// 房间是否存在
func (impl *ModuleImpl) IsRoomExist(roomID Define.RoomID) bool {
	_, ok := impl.id2rooms[roomID]
	return ok
}

// 房间是否包含该名称玩家
func (impl *ModuleImpl) IsRoomContainsUser(roomID Define.RoomID, user Define.UserName) bool {
	room := impl.id2rooms[roomID]
	for _, u := range room.Users {
		if u.Name == user {
			return true
		}
	}
	return false
}

const (
	ERoomNotExist      = 1 // 房间名不存在
	EUserNotMatchRoom  = 2 // 用户不在该目标房间中
	ESameNameParameter = 3 // 已经存在同名目标 （房间/用户）

	EChangeReadyRetChangeFailed = 11 // 无法在此时更改准备状态，更改失败
)

func (impl *ModuleImpl) New() *ModuleImpl {
	impl.id2rooms = make(map[Define.RoomID]*Define.Room)
	return impl
}

func (impl *ModuleImpl) ConnectToSvr(svr *Define.Server) {
	impl.svr = svr
}

// 开始
func (impl *ModuleImpl) Start() {

}

// 创建房间
func (impl *ModuleImpl) Create(roomID Define.RoomID, firstUser Define.UserName) int {
	if impl.IsRoomExist(roomID) {
		// 房间名相同错误
		return ESameNameParameter
	}

	room := Define.Room{
		RoomID: roomID,
		Users: []*Define.User{
			&Define.User{Name: firstUser},
		},
	}
	impl.id2rooms[roomID] = &room
	return 0
}

// 玩家加入房间
func (impl *ModuleImpl) Join(roomID Define.RoomID, username Define.UserName) int {
	if !impl.IsRoomExist(roomID) {
		return ERoomNotExist
	}

	if impl.IsRoomContainsUser(roomID, username) {
		// 用户名相同加入失败
		return ESameNameParameter
	}

	room := impl.id2rooms[roomID]
	room.Users = append(room.Users, &Define.User{
		Name: username,
	})

	impl.BroadCast(roomID, Define.WSCRoomUserJoin, map[string]string{
		"username": string(username),
	}, username)

	return 0
}

// 玩家离开房间
func (impl *ModuleImpl) Leave(roomID Define.RoomID, user Define.UserName) int {
	if !impl.IsRoomExist(roomID) {
		return ERoomNotExist
	}

	if !impl.IsRoomContainsUser(roomID, user) {
		return EUserNotMatchRoom
	}

	room := impl.id2rooms[roomID]
	for index, u := range room.Users {
		if u.Name == user {
			room.Users = append(room.Users[:index], room.Users[index+1:]...)
		}
	}

	var pushModule, ok = impl.svr.GetModule(Define.ModuleNamePush).(Define.PushModule)
	if ok {
		pushModule.RemoveNotify(roomID, user)
	}

	impl.BroadCast(roomID, Define.WSCRoomUserLeave, map[string]string{
		"username": string(user),
	}, user)

	// 玩家已全部离开，删除房间
	if len(room.Users) == 0 {
		delete(impl.id2rooms, roomID)
	}
	return 0
}

// 玩家准备
func (impl *ModuleImpl) ChangeReady(roomID Define.RoomID, user Define.UserName, ready bool) int {
	if !impl.IsRoomExist(roomID) {
		return ERoomNotExist
	}

	if !impl.IsRoomContainsUser(roomID, user) {
		return EUserNotMatchRoom
	}

	room := impl.id2rooms[roomID]
	// 已经开始游戏了，禁止再更改状态
	if room.RoomState == Define.RoomStateInGame {
		return EChangeReadyRetChangeFailed
	}

	var allReady = true
	for _, u := range room.Users {
		if u.Name == user {
			u.Ready = ready
			if !ready {
				allReady = false
				break
			}
		}

		if !u.Ready {
			allReady = false
		}
	}

	impl.BroadCast(roomID, Define.WSCRoomUserChangeReady, map[string]string{
		"username": string(user),
		"ready":    strconv.FormatBool(ready),
	}, user)

	if allReady {
		// 所有玩家准备完成，开始游戏
		impl.BroadCastAll(roomID, Define.WSCRoomUserChangeReady, map[string]string{

		})
	}

	return 0
}

// 取房间所有玩家信息
func (impl *ModuleImpl) GetUserList(roomID Define.RoomID) (*[]Define.User, int) {
	if !impl.IsRoomExist(roomID) {
		return nil, ERoomNotExist
	}

	room := impl.id2rooms[roomID]
	var users = make([]Define.User, len(room.Users))
	for i, u := range room.Users {
		users[i] = *u
	}
	return &users, 0

}

func (impl *ModuleImpl) BroadCastAll(roomID Define.RoomID, code int, param interface{}) {
	impl.BroadCast(roomID, code, param, "")
}

func (impl *ModuleImpl) BroadCast(roomID Define.RoomID, code int, param interface{}, ignoreName Define.UserName) {
	var pushModule, ok = impl.svr.GetModule(Define.ModuleNamePush).(Define.PushModule)
	if ok {
		pushModule.Broadcast(&Define.BroadcastPacket{
			RoomID: roomID,
			Packet: Define.WSPacket{
				Type:     Define.WSPacketTypeRawMsg,
				Category: Define.WSPacketCategoryRoom,
				Code:     code,
				Param:    param,
			},
			IgnoreUser: ignoreName,
		})
	}
}
