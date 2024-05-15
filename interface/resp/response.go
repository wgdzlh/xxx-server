package resp

import (
	"net/http"

	"xxx-server/domain/entity"

	"github.com/gin-gonic/gin"
)

const (
	// Success Success
	Success = 200
	// ErrResponse Error
	ErrResponse = 500
	// ErrBadRequest 请求参数错误
	ErrBadRequest = 400
	// ErrUnauthorized 非法请求
	ErrUnauthorized = 401
	// ErrInternalServer 服务端错误
	// ErrInternalServer = 501
)

var (
	emptyObject = struct{}{}

	code2Msg = map[int]string{
		Success:         "ok",
		ErrResponse:     "error",
		ErrBadRequest:   "bad request",
		ErrUnauthorized: "unauthorized",
	}
)

type Response struct {
	Code int    `json:"code" enums:"200,400,401,500"` // code
	Msg  string `json:"msg"`                          // message
	Data any    `json:"data" swaggertype:"object"`    // data
}

type IdData struct {
	Id uint64 `json:"id"`
}

type NumData struct {
	Num int64 `json:"num"`
}

type ListData struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
}

// below response structs are only for swagger doc
type IdResponse struct {
	Response
	Data IdData `json:"data"`
}

type BgTaskResponse struct {
	Response
	Data entity.BackgroundTask `json:"data"`
}

type BgTasksResponse struct {
	Response
	Data struct {
		List  []entity.BackgroundTask `json:"list"`
		Total int64                   `json:"total"`
	} `json:"data"`
}

type CronTaskResponse struct {
	Response
	Data entity.CronTask `json:"data"`
}

type CronTasksResponse struct {
	Response
	Data struct {
		List  []entity.CronTask `json:"list"`
		Total int64             `json:"total"`
	} `json:"data"`
}

type SettingResponse struct {
	Response
	Data entity.Setting `json:"data"`
}

type SettingsResponse struct {
	Response
	Data []entity.Setting `json:"data"`
}

// get msg by code
func GetMsg(code int) string {
	if msg, ok := code2Msg[code]; ok {
		return msg
	}
	return "unknown"
}

// gin resp to json
func JSON(c *gin.Context, code int, data ...any) {
	resp := Response{
		Code: code,
		Msg:  GetMsg(code),
	}
	if len(data) == 1 {
		resp.Data = data[0]
	}
	// 保证返回时data字段不为nil
	if resp.Data == nil {
		resp.Data = emptyObject
	}
	c.JSON(http.StatusOK, resp)
}

// gin resp to json with msg
func JSONWithMsg(c *gin.Context, code int, msg string, data ...any) {
	resp := Response{
		Code: code,
		Msg:  msg,
	}
	if len(data) == 1 {
		resp.Data = data[0]
	}
	// 保证返回时data字段不为nil
	if resp.Data == nil {
		resp.Data = emptyObject
	}
	c.JSON(http.StatusOK, resp)
}

func AbortJSONWithMsg(c *gin.Context, code int, msg string, data ...any) {
	resp := Response{
		Code: code,
		Msg:  msg,
	}
	if len(data) == 1 {
		resp.Data = data[0]
	}
	// 保证返回时data字段不为nil
	if resp.Data == nil {
		resp.Data = emptyObject
	}
	c.AbortWithStatusJSON(http.StatusOK, resp)
}
