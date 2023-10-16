package tool

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// 生成范围随机float64
func RandFloat64(min, max float64) float64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	minStr := strconv.FormatFloat(min, 'f', -1, 64)
	// 不包含小数点
	if strings.Index(minStr, ".") == -1 {
		return max
	}
	multipleNum := len(minStr) - (strings.Index(minStr, ".") + 1)
	multiple := math.Pow10(multipleNum)
	minMult := min * multiple
	maxMult := max * multiple
	randVal := RandInt64(int64(minMult), int64(maxMult))
	result := float64(randVal) / multiple
	return result
}

// 随机整数
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min+1) + min
}
