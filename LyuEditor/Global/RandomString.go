package Global

import (
	"math/rand"
	"time"
)

func RandStringRunes() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, letterLenth)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}