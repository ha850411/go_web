package service

import (
	"time"

	"github.com/gin-gonic/gin"
)

func GetCommonOutput(c *gin.Context, active string) map[string]interface{} {
	output := make(map[string]interface{})
	username, _ := c.Get("username")
	output["username"] = username
	output["active"] = active
	output["staticFreshFlag"] = time.Now().Unix()
	output["menu"] = GetMenu()
	return output
}
