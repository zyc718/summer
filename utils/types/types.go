package types

import "time"

// TimeNow returns current wall time in UTC rounded to milliseconds.
//返回毫秒
func TimeNow() time.Time {
	//	设置时区
	l, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().UTC().Round(time.Millisecond).In(l)
}
