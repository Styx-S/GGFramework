package Room

import (
	"GGFramework/Define"
	"github.com/gin-gonic/gin"
	"net/http"
)


type JoinRequest struct {
	RoomID string `form:"room_id"`
	Username string `form:"username"`
	Create bool `form:"create"`
}
const (
	JoinRoomErrorRoomNotExist     = 1
	JoinRoomErrorCreateRoomFailed = 2
	JoinRoomErrorHasSameUserName  = 3
)
// 玩家加入房间
func JoinRoom(ctx *gin.Context) {
	var request JoinRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.JSON(http.StatusOK, Define.Response{
			Ret: Define.ParametersError,
			Msg: err.Error(),
		})
		return
	}

	ret := Join(ID(request.RoomID), UserName(request.Username))
	if ret == JoinRetSucceed {
		ctx.JSON(http.StatusOK, Define.Response{
			Ret: Define.NoError,
			Msg: "加入房间成功",
		})
		return
	}
	if ret == JoinRetNotExist {
		if request.Create {
			createRet := Create(ID(request.RoomID), UserName(request.Username))
			if createRet == CreateRetSucceed {
				ctx.JSON(http.StatusOK, Define.Response{
					Ret: Define.NoError,
					Msg: "创建房间成功",
				})
				return
			} else {
				ctx.JSON(http.StatusOK, Define.Response{
					Ret: JoinRoomErrorCreateRoomFailed,
					Msg: "创建房间失败!",
				})
				return
			}
		} else {
			ctx.JSON(http.StatusOK, Define.Response{
				Ret: JoinRoomErrorRoomNotExist,
				Msg: "该房间不存在，无法加入!",
			})
			return
		}
	}
	if ret == JoinRetSameName {
		ctx.JSON(http.StatusOK, Define.Response{
			Ret: JoinRoomErrorHasSameUserName,
			Msg: "房间中已经有了同名玩家，加入失败!",
		})
		return
	}

	ctx.JSON(http.StatusOK, Define.Response{
		Ret: Define.UndefinedError,
	})
	return
}

type LeaveRequest struct {
	RoomID string
	Username string
}
// 玩家离开房间
func LeaveRoom(ctx *gin.Context) {
	var request LeaveRequest
	err := ctx.ShouldBind(&request)
	if err != nil {

	}
}