package tidbTools

import (
	"fmt"
	"strconv"
	"strings"
)

func ToFloat64(str string) float64 {
	var newNum float64
	tmpStr := strings.ToLower(str)
	if strings.Contains(tmpStr, "e+") || strings.Contains(tmpStr, "e-") {
		_, err := fmt.Sscanf(tmpStr, "%e", &newNum)
		if err != nil {
			fmt.Printf("fmt.Sscanf error, numStr:%s, err:%v", tmpStr, err)
			return 0
		}
	} else {
		newNum, _ = strconv.ParseFloat(str, 64)
	}
	return newNum
}

func Add(a, b int) int {
	return a + b
}
