package v1

import (
	"crypto/md5"

	"github.com/gin-gonic/gin"
)

// AlphanumericSet 字母数字集
var AlphanumericSet = []rune{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
}

// GetInvCodeByUID 获取指定长度的邀请码
func GetInvCodeByUID(uid string, l int) string {
	// 因为 md5 值为 16 字节
	if l > 16 {
		return ""
	}
	sum := md5.Sum([]byte(uid))
	var code []rune
	for i := 0; i < l; i++ {
		idx := sum[i] % byte(len(AlphanumericSet))
		code = append(code, AlphanumericSet[idx])
	}
	return string(code)
}

func GetInviteCode(c *gin.Context) {
	address := c.Params.ByName("code")
	// TODO check address

	code := GetInvCodeByUID(address, 6)

	c.JSON(200, gin.H{
		"message": "获取",
		"code":    200,
		"data":    code,
	})
	return
}
