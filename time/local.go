package time

import (
	"github.com/danakum/go-util/config"
	"time"
)

func ToLocal(t time.Time) time.Time {
	return t.In(config.AppConf.Location)
}
