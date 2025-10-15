package utils

import (
	"9zhou-scripts/client"
	"time"
)

// ParseTime 解析多种时间格式
func ParseTime(timeStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.ParseInLocation(format, timeStr, time.Local); err == nil {
			return t, nil
		}
	}

	return time.Time{}, &time.ParseError{Layout: "", Value: timeStr, LayoutElem: "", ValueElem: timeStr, Message: "cannot parse time"}
}

// parseTime 解析多种时间格式 (为了兼容已有的代码)
func parseTime(timeStr string) (time.Time, error) {
	return ParseTime(timeStr)
}

func SumAmount(order *client.GetHXOrderResp) float64 {
	amount := 0.0
	for i := range order.Data {
		amount += float64(order.Data[i].TotalPrice)
	}
	return amount
}
