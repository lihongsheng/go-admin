// Package response 统一响应格式
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

const (
	CodeOK            = 0
	CodeError         = 7
	CodeNotInstalled  = 9001
	CodeAlreadyDone   = 9002
	CodeUnauthorized  = 401
	CodeForbidden     = 403
)

func OK(c *gin.Context, data interface{})            { c.JSON(http.StatusOK, Body{CodeOK, "ok", data}) }
func OKMsg(c *gin.Context, msg string)               { c.JSON(http.StatusOK, Body{CodeOK, msg, nil}) }
func Fail(c *gin.Context, msg string)                { c.JSON(http.StatusOK, Body{CodeError, msg, nil}) }
func FailCode(c *gin.Context, code int, msg string)  { c.JSON(http.StatusOK, Body{code, msg, nil}) }
func FailHTTP(c *gin.Context, status, code int, msg string) {
	c.JSON(status, Body{code, msg, nil})
}
