package http_handler

import (
	"net/http"
	log "xxx-server/application/logger"
	"xxx-server/interface/resp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	MULTIPART_FILE_KEY = "file"
)

// func getUserId(c *gin.Context) int64 {
// 	return int64(c.GetInt(middleware.USER_ID_CTX))
// }

func getQueryFilter(c *gin.Context) map[string]string {
	var filter = map[string]string{}
	queries := c.Request.URL.Query()
	for key, value := range queries {
		if value[0] != "" {
			filter[key] = value[0]
		}
	}
	return filter
}

func errIdResp(c *gin.Context, tag, idStr string) {
	msg := tag + ":id不合法"
	log.Error(msg, zap.String("id", idStr))
	resp.JSONWithMsg(c, resp.ErrBadRequest, msg)
}

func errParamResp(c *gin.Context, tag string, err error) {
	msg := tag + ":参数解析失败:"
	log.Error(msg, zap.Error(err))
	resp.JSONWithMsg(c, resp.ErrBadRequest, msg+err.Error())
}

// func errInvalidResp(c *gin.Context, tag string, err error) {
// 	msg := tag + ":参数不合法:"
// 	log.Error(msg, zap.Error(err))
// 	resp.JSONWithMsg(c, resp.ErrBadRequest, msg+err.Error())
// }

func errOpResp(c *gin.Context, tag string, err error) {
	msg := tag + ":操作失败:"
	log.Error(msg, zap.Error(err))
	resp.JSONWithMsg(c, resp.ErrResponse, msg+err.Error())
}

func errIdRespStatus(c *gin.Context, tag, idStr string) {
	msg := tag + ":id不合法"
	log.Error(msg, zap.String("id", idStr))
	c.AbortWithStatus(http.StatusBadRequest)
}

func errOpRespStatus(c *gin.Context, tag string, err error) {
	msg := tag + ":操作失败"
	log.Error(msg, zap.Error(err))
	c.AbortWithStatus(http.StatusInternalServerError)
}

func errNotFoundRespStatus(c *gin.Context, tag string, err error) {
	msg := tag + ":资源不存在"
	log.Error(msg, zap.Error(err))
	c.AbortWithStatus(http.StatusNotFound)
}
