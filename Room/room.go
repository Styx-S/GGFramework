package Room

type ID string
type UserName string

type User struct {
	Name UserName
	Ready bool
}

type Room struct {
	RoomID ID
	Users []UserName
}


var id2rooms = make(map[ID]Room)

// 房间是否存在
func IsRoomExist(roomID ID) bool {
	_, ok := id2rooms[roomID]
	return ok
}

// 房间是否包含该名称玩家
func IsRoomContainsUser(roomID ID, user UserName) bool {
	room := id2rooms[roomID]
	for _, u := range room.Users{
		if u == user {
			return true
		}
	}
	return false
}

type CreateRet int

const (
	CreateRetSucceed = 0
	CreateRetSameName = 1
)
// 创建房间
func Create(roomID ID, firstUser UserName) CreateRet {
	if IsRoomExist(roomID) {
		return CreateRetSameName
	}


	room := Room{
		RoomID: roomID,
		Users: []UserName{firstUser},
	}
	id2rooms[roomID] = room
	return CreateRetSucceed
}


type JoinRet int
const (
	JoinRetSucceed = 0
	JoinRetNotExist = 1
	JoinRetSameName = 2

)
// 玩家加入房间
func Join(roomID ID, user UserName) JoinRet {
	if !IsRoomExist(roomID) {
		return JoinRetNotExist
	}

	if IsRoomContainsUser(roomID, user) {
		return JoinRetSameName
	}


	room := id2rooms[roomID]
	room.Users = append(room.Users, user)
	// 不加这行不行，就不能是引用吗？我再学下go
	id2rooms[roomID] = room
	return JoinRetSucceed
}

type LeveRet int
const (
	LeveRetSucceed = 0
	LeveRetRoomNotExist = 1
	LeveRetUserNotExist = 2

)
func Leave(roomID ID, user UserName) {

}
