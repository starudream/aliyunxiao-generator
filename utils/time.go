package utils

import (
	"time"

	"github.com/golang-module/carbon/v2"
)

type Carbon = carbon.Carbon

func NewTime(ts ...time.Time) Carbon {
	tt := time.Now()
	if len(ts) > 0 {
		tt = ts[0]
	}
	return carbon.CreateFromStdTime(tt).SetTimezone(carbon.PRC).SetLocale("zh-CN").SetWeekStartsAt(carbon.Monday)
}

func ParseTime(value, layout string) (Carbon, error) {
	t := NewTime().ParseByLayout(value, layout)
	return t, t.Error
}
