package parameter

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetParameterValueUint(c *gin.Context, key string) (uint, error) {
	value, err := strconv.ParseUint(c.Param(key), 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid %s", key)
	}
	return uint(value), nil
}
