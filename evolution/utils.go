package evolution

import (
	"math/rand"
	"sync"
	"time"
)

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

//func RandString(n int) string {
//	x := sync.Mutex{}
//	x.Lock()
//
//	sb := strings.Builder{}
//	sb.Grow(n)
//	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
//
//	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
//		if remain == 0 {
//			cache, remain = src.Int63(), letterIdxMax
//		}
//		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
//			sb.WriteByte(letterBytes[idx])
//			i--
//		}
//		cache >>= letterIdxBits
//		remain--
//	}
//	s := sb.String()
//	x.Unlock()
//
//	return s
//}

func RandString(n int) string {

	b := make([]byte, n)
	for i := range b {
		x := sync.Mutex{}
		x.Lock()
		int63 := rand.Int63()
		x.Unlock()
		b[i] = letterBytes[int63%int64(len(letterBytes))]
	}
	return string(b)
}
