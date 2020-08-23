package models

import (
	"strings"
	"time"
)

var CustomTimeFormatter = "01-02-2006 15:04:05 Mon"

type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalJSON(buf []byte) error {
	tt, err := time.Parse(time.RubyDate, strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

func (t CustomTime) Format() string {
	return time.Time(t.Time).Format(CustomTimeFormatter)
}

type Tweet struct {
	CreatedAt CustomTime `json:"created_at"`
	Text      string     `json:"text"`
}
