package utils

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// JsonpToJson modify jsonp string to json string
// Example: forbar({a:"1",b:2}) to {"a":"1","b":2}
func JsonpToJson(json string) string {
	start := strings.Index(json, "{")
	end := strings.LastIndex(json, "}")
	start1 := strings.Index(json, "[")
	if start1 > 0 && start > start1 {
		start = start1
		end = strings.LastIndex(json, "]")
	}
	if end > start && end != -1 && start != -1 {
		json = json[start : end+1]
	}
	json = strings.Replace(json, "\\'", "", -1)
	regDetail, _ := regexp.Compile("([^\\s\\:\\{\\,\\d\"]+|[a-z][a-z\\d]*)\\s*\\:")
	return regDetail.ReplaceAllString(json, "\"$1\":")
}

func ToString(v interface{}) string {
	switch f := v.(type) {
	case bool:
		if f {
			return "true"
		} else {
			return "false"
		}
	case float32:
		return strconv.FormatFloat(float64(f), 'E', -1, 32)
	case float64:
		return strconv.FormatFloat(f, 'E', -1, 64)
	case int:
		return strconv.Itoa(f)
	case int8:
		return strconv.FormatInt(int64(f), 10)
	case int16:
		return strconv.FormatInt(int64(f), 10)
	case int32:
		return strconv.FormatInt(int64(f), 10)
	case int64:
		return strconv.FormatInt(f, 10)
	case uint:
		return strconv.FormatUint(uint64(f), 10)
	case uint8:
		return strconv.FormatUint(uint64(f), 10)
	case uint16:
		return strconv.FormatUint(uint64(f), 10)
	case uint32:
		return strconv.FormatUint(uint64(f), 10)
	case uint64:
		return strconv.FormatUint(f, 10)
	case time.Time:
		return f.Format("2006-01-02 15:04:05")
	case string:
		return f
	default:
		return ""
	}
}

func BytesToSize(length int64) string {
	var k = 1024 // or 1024
	var sizes = []string{"Bytes", "KB", "MB", "GB", "TB"}
	if length == 0 {
		return "0 Bytes"
	}
	i := math.Floor(math.Log(float64(length)) / math.Log(float64(k)))
	r := float64(length) / math.Pow(float64(k), i)
	return strconv.FormatFloat(r, 'f', 3, 64) + " " + sizes[int(i)]
}
