package Global

import (
	"math/rand"
	"time"
)

// 生成随机字符，作为html命名
func RandStringRunes() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, letterLenth)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}