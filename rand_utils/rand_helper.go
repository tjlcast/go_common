package rand_utils

import (
	"math/rand"
	"time"
)

func GenRandInt(max int) int {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(max) // 生成0~9的随机
	return randomNum
}
