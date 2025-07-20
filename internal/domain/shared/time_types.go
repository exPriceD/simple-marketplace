package shared

import "time"

type Timestamp = time.Time

func NowUTC() time.Time {
	return time.Now().UTC()
}
