package timer

import (
	"strconv"
	"time"
)

// 获取当前时间毫秒级和13位的时间戳,UTC ,数组截取时为了避免尾数四舍五入的情况
func GetCurTime() (now time.Time, ts int64) {
	ts, _ = strconv.ParseInt(strconv.FormatInt(time.Now().UnixNano(),
		10)[:13], 10, 64)
	now = time.Unix(0, ts*1000000).UTC()
	return
}

func TimeToTimeStampStr(t time.Time) string {
	return strconv.FormatInt(t.UTC().UnixNano(),
		10)[:13]
}

func TimeToTimeStampInt64(t time.Time) int64 {
	ts, _ := strconv.ParseInt(strconv.FormatInt(t.UTC().UnixNano(),
		10)[:13], 10, 64)
	return ts
}

func TimeStampToTime(ts int64) time.Time {
	return time.Unix(0, ts*1000000).UTC()
}
