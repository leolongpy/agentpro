package utils

import (
	"fmt"
)

func PK(endpoint, metric string, tags map[string]string) string {
	if tags == nil || len(tags) == 0 {
		return fmt.Sprintf("%s/%s", endpoint, metric)
	}
	return fmt.Sprintf("%s/%s/%s", endpoint, metric, SortedTags(tags))
}

func UUID(endpoint, metric string, tags map[string]string, dstype string, step int) string {
	if tags == nil || len(tags) == 0 {
		return fmt.Sprintf("%s/%s/%s/%d", endpoint, metric, dstype, step)
	}
	return fmt.Sprintf("%s/%s/%s/%s/%d", endpoint, metric, SortedTags(tags), dstype, step)
}

func Checksum(endpoint string, metric string, tags map[string]string) string {
	pk := PK(endpoint, metric, tags)
	return Md5(pk)
}

func Strip(s_ string, chars_ string) string {
	s, chars := []rune(s_), []rune(chars_)
	length := len(s)
	max := len(s) - 1
	l, r := true, true
	start, end := 0, max
	tmpEnd := 0
	charset := make(map[rune]bool)
	for i := 0; i < len(chars); i++ {
		charset[chars[i]] = true
	}
	for i := 0; i < length; i++ {
		if _, exist := charset[s[i]]; l && !exist {
			start = i
			l = false
		}
		tmpEnd = max - i
		if _, exist := charset[s[tmpEnd]]; r && !exist {
			end = tmpEnd
			r = false
		}
		if !l && !r {
			break
		}
	}
	if l && r {
		return ""
	}
	return string(s[start : end+1])
}
