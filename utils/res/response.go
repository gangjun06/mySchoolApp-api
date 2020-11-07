package res

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrType int

const (
	ERR_BAD_REQUEST ErrType = 1 + iota
	ERR_SERVER
	ERR_DUPLICATE_NAME
	ERR_DUPLICATE_EMAIL
	ERR_AUTH
	ERR_EXPIRED
	ERR_ALREADY_VERIFIED
	ERR_NOT_MATCH
	ERR_NOT_FOUND
	ERR_NO_PERMISSION
)

func (r *res) SendError(errType ErrType, text string) {
	var ErrCode string
	var Status int

	set := func(errCode string, status int) {
		ErrCode = errCode
		Status = status
	}

	switch errType {
	case ERR_BAD_REQUEST:
		set("ERR_BAD_REQUEST", http.StatusBadRequest)
	case ERR_SERVER:
		set("ERR_SERVER", http.StatusInternalServerError)
	case ERR_DUPLICATE_NAME:
		set("ERR_DUPLICATE_NAME", http.StatusConflict)
	case ERR_DUPLICATE_EMAIL:
		set("ERR_DUPLICATE_EMAIL", http.StatusConflict)
	case ERR_AUTH:
		set("ERR_AUTH", http.StatusUnauthorized)
	case ERR_EXPIRED:
		set("ERR_EXPIRED", http.StatusBadRequest)
	case ERR_ALREADY_VERIFIED:
		set("ERR_ALREADY_VERIFIED", http.StatusBadRequest)
	case ERR_NOT_MATCH:
		set("ERR_NOT_MATCH", http.StatusBadRequest)
	case ERR_NOT_FOUND:
		set("ERR_NOT_FOUND", http.StatusBadRequest)
	case ERR_NO_PERMISSION:
		set("ERR_NO_PERMISSION", http.StatusUnauthorized)
	}

	r.c.JSON(Status, gin.H{
		"code":    ErrCode,
		"message": text,
	})
	r.c.Abort()
}
