package _func

import (
	"sort"
	"strconv"
	"strings"
)

func ArrayFilter(arr []string) ([]string) {
	//result := []string{}
	// 必须确保数组从小到大。
	sort.Strings(arr)
	result2 := []string{}
	tempMap := map[string]int{}
	for _, e := range arr {
		//k := len(tempMap)
		tempMap[e] += 1
		//if len(tempMap) != k {
		//	result = append(result, e)
		//}
	}
	for w, f := range tempMap {
		k := []string{w, strconv.Itoa(f)}
		b := strings.Join(k, "(")
		b = b + ")"
		result2 = append(result2, b)
	}
	return result2
}