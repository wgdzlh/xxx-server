package dto

import (
	"github.com/gin-gonic/gin"
	json "github.com/json-iterator/go"
)

func BindJsonAndToBs(c *gin.Context, input any) (out []byte, err error) {
	if err = c.ShouldBindJSON(input); err != nil {
		return
	}
	out, err = json.Marshal(input)
	return
}

func BindJsonAndToKeyedBs(c *gin.Context, input any, key string) (out []byte, err error) {
	if err = c.ShouldBindJSON(input); err != nil {
		return
	}
	out, err = ToKeyedBs(input, key)
	return
}

func ToKeyedBs(input any, key string) (out []byte, err error) {
	out, err = json.Marshal(struct {
		Key   string `json:"key"`
		Value any    `json:"value"`
	}{key, input})
	return
}
