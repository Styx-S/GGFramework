package Room

import (
	"GGFramework/Define"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	ErrorRoomNotExist      = 101
	ErrorCreateRoomFailed  = 102
	ErrorHasSameUserName   = 103
	ErrorUserNotInRoom     = 104
	ErrorChangeStateFailed = 105
)

type JoinRequest struct {
	RoomID   string `form:"room_id"`
	Username string `form:"username"`
	Create   bool   `form:"create"`
}

var ModuleInstance Define.RoomModule = (&ModuleImpl{}).New()

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

	ret := ModuleInstance.Join(Define.RoomID(request.RoomID), Define.UserName(request.Username))
	if ret == 0 {
		// 查当前房间内所有人信息
		userInfos, _ := ModuleInstance.GetUserList(Define.RoomID(request.RoomID))

		ctx.JSON(http.StatusOK, Define.Response{
			Ret:   Define.NoError,
			Msg:   "加入房间成功",
			Param: convUserRspList(userInfos),
		})
		return
	}
	if ret == ERoomNotExist {
		if request.Create {
			createRet := ModuleInstance.Create(Define.RoomID(request.RoomID), Define.UserName(request.Username))
			if createRet == 0 {
				ctx.JSON(http.StatusOK, Define.Response{
					Ret: Define.NoError,
					Msg: "创建房间成功",
				})
				return
			} else {
				ctx.JSON(http.StatusOK, Define.Response{
					Ret: ErrorCreateRoomFailed,
					Msg: "创建房间失败!",
				})
				return
			}
		} else {
			ctx.JSON(http.StatusOK, Define.Response{
				Ret: ErrorRoomNotExist,
				Msg: "该房间不存在，无法加入!",
			})
			return
		}
	}
	if ret == ESameNameParameter {
		ctx.JSON(http.StatusOK, Define.Response{
			Ret: ErrorHasSameUserName,
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
	RoomID   string `form:"room_id"`
	Username string `form:"username"`
}

// 玩家离开房间
func LeaveRoom(ctx *gin.Context) {
	var request LeaveRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.JSON(http.StatusOK, Define.Response{
			Ret: Define.ParametersError,
			Msg: err.Error(),
		})
		return
	}

	ret := ModuleInstance.Leave(Define.RoomID(request.RoomID), Define.UserName(request.Username))
	switch ret {
	case 0:
		ctx.JSON(http.StatusOK, Define.Response{
			Ret: Define.NoError,
			Msg: "离开房间成功",
		})
		return
	case ERoomNotExist:
		ctx.JSON(http.StatusOK, Define.Response{
			Ret: ErrorRoomNotExist,
			Msg: "离开房间失败，该房间不存在！",
		})
		return
	case EUserNotMatchRoom:
		ctx.JSON(http.StatusOK, Define.Response{
			Ret: ErrorUserNotInRoom,
			Msg: "离开房间失败，该用户不在该房间中！",
		})
		return
	}
}

type ReadyRequest struct {
	RoomID   string `form:"room_id"`
	Username string `form:"username"`
	Ready    bool   `form:"ready"`
}

const ()

// 更改准备状态
func ReadyRoom(ctx *gin.Context) {
	var request ReadyRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		ctx.JSON(http.StatusOK, Define.Response{
			Ret: Define.ParametersError,
			Msg: err.Error(),
		})
		return
	}

	ret := ModuleInstance.ChangeReady(Define.RoomID(request.RoomID), Define.UserName(request.Username), request.Ready)
	switch ret {
	case 0:
		ctx.JSON(http.StatusOK, Define.Response{
			Ret: Define.NoError,
			Msg: "更改状态成功",
			Param: map[string]string{
				"ready": strconv.FormatBool(request.Ready),
			},
		})
		return
	case EChangeReadyRetChangeFailed:
		ctx.JSON(http.StatusOK, Define.Response{
			Ret: ErrorChangeStateFailed,
			Msg: "修改状态失败！",
		})
		return
	case ERoomNotExist:
		ctx.JSON(http.StatusOK, Define.Response{
			Ret: ErrorRoomNotExist,
			Msg: "修改状态失败，房间不存在！",
		})
		return
	case EUserNotMatchRoom:
		ctx.JSON(http.StatusOK, Define.Response{
			Ret: ErrorUserNotInRoom,
			Msg: "修改状态失败，该用户不在该房间中！",
		})
		return
	}
}

/* 工具方法 */
func convUserRsp(user *Define.User) Define.UserRsp {
	return Define.UserRsp{
		Username:   string(user.Name),
		ReadyState: user.Ready,
	}
}

func convUserRspList(users *[]Define.User) []Define.UserRsp {
	rps := make([]Define.UserRsp, len(*users))
	for i, u := range *users {
		rps[i] = convUserRsp(&u)
	}
	return rps
}
