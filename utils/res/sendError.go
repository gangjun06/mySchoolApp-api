package res

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type R map[string]interface{}

type res struct {
	c *gin.Context
}

func New(c *gin.Context) *res {
	return &res{c}
}

type SuccessType int

const (
	SUCCESS_OK SuccessType = 1 + iota
	SUCCESS_ACCEPTED
	SUCCESS_CREATED
	SUCCESS_NO_CONTENT
)

func (r *res) Response(successType SuccessType, data interface{}) {

	var SuccessCode string
	var Status int

	set := func(errCode string, status int) {
		SuccessCode = errCode
		Status = status
	}

	switch successType {
	case SUCCESS_OK:
		set("SUCCESS_OK", http.StatusOK)
	case SUCCESS_CREATED:
		set("SUCCESS_CREATED", http.StatusCreated)
	case SUCCESS_ACCEPTED:
		set("SUCCESS_ACCEPTED", http.StatusAccepted)
	case SUCCESS_NO_CONTENT:
		r.c.JSON(http.StatusNoContent, R{})
	}

	m := make(map[string]interface{})
	m["code"] = SuccessCode
	m["message"] = ""
	j, _ := json.Marshal(data)
	json.Unmarshal(j, &m)
	r.c.JSON(Status, m)
}
