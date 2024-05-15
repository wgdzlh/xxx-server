package entity

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"xxx-server/application/utils"
	"xxx-server/infrastructure/config"
)

const (
	emptyStr = `""`
)

var (
	EmptyTime  = JSONTime{}
	emptyStrBs = []byte(emptyStr)

	ErrUnknownSqlValueType = errors.New("unknown sql result value type")
)

type JSONTime struct {
	time.Time
}

func (jt JSONTime) MarshalJSON() ([]byte, error) {
	if jt == EmptyTime {
		return emptyStrBs, nil
	}
	return fmt.Appendf(nil, `"%s"`, utils.TimeToStr(jt.Time)), nil
}

func (jt *JSONTime) UnmarshalJSON(bs []byte) error {
	if jt == nil {
		return nil
	}
	if string(bs) == emptyStr {
		return nil
	}
	_t, err := utils.StrToTime(utils.B2S(bytes.Trim(bs, `"`)))
	if err == nil {
		*jt = JSONTime{_t}
	}
	return err
}

func (jt JSONTime) Value() (driver.Value, error) {
	return jt.Time, nil
}

func (jt *JSONTime) Scan(value any) (err error) {
	if value == nil {
		jt.Time = time.Time{}
		return
	}
	switch v := value.(type) {
	case time.Time:
		jt.Time = v
	case string:
		jt.Time, err = utils.StrToTime(v)
	case []byte:
		jt.Time, err = utils.StrToTime(utils.B2S(v))
	default:
		err = ErrUnknownSqlValueType
	}
	return
}

func (jt JSONTime) ToDateStr() string {
	// if jt == nil {
	// 	return ""
	// }
	return jt.Format(config.Q_DATE_FORMAT)
}

func (jt JSONTime) ToDateTimeStr() string {
	return jt.Format(config.Q_TIME_FORMAT)
}

func (jt JSONTime) ToDateAndTime() (string, int) {
	return jt.Format(config.Q_DATE_FORMAT), jt.Hour()*100 + jt.Minute()
}

func (jt *JSONTime) FromStr(ts, format string) (err error) {
	if jt == nil {
		return
	}
	jt.Time, err = utils.CstStrToTime(ts, format)
	return
}

func StrToJT(ts, format string) (jt JSONTime) {
	jt.Time, _ = utils.CstStrToTime(ts, format)
	return
}

type Date string

func (d *Date) Scan(value any) error {
	if value == nil {
		*d = ""
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*d = Date(v.Format(config.Q_DATE_FORMAT))
	case string:
		*d = Date(v)
	}
	return nil
}
