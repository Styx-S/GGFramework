package Define

/* room */

type RoomID string
type UserName string

type User struct {
	Name  UserName
	Ready bool
}

type RoomState int

const (
	RoomStateNormal = 1
	RoomStateInGame = 2
)

type Room struct {
	RoomID    RoomID
	Users     []*User
	RoomState RoomState
}

/* push */

// see WSPacket (response.go)

type BroadcastPacket struct {
	RoomID     RoomID
	Packet     WSPacket
	IgnoreUser UserName
}
