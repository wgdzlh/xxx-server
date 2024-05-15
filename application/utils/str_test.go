package utils

import (
	"strings"
	"testing"

	json "github.com/json-iterator/go"
)

func TestXxx(t *testing.T) {
	s := "站点ID"
	g, e := Utf8StrToGbk(s)
	t.Log(g, e)
	g, e = GbkStrToUtf8(g)
	t.Log(g, e)
	t.Log(Utf8StrToGbk(s))
	s = "区域"
	g, e = Utf8StrToGbk(s)
	t.Log(g, e)
	g, e = GbkStrToUtf8(g)
	t.Log(g, e)
	t.Log(StrToUint64(""))
}

func TestDegree(t *testing.T) {
	s := `23°26' 22" N`
	t.Log(DegreeSecToFloat(s))
}

func TestDateStr(t *testing.T) {
	s := "2020"
	t.Log(DateStrToTime(s))
	s = "202011"
	t.Log(DateStrToTime(s))
	s = "20201113"
	t.Log(DateStrToTime(s))
	s = "2020Q4"
	t.Log(DateStrToTime(s))
	s = "2020Q5"
	t.Log(DateStrToTime(s))
	s = "201"
	t.Log(DateStrToTime(s))
	s = "2020-01"
	t.Log(DateStrToTime(s))
	s = "2020A1"
	t.Log(DateStrToTime(s))
	s = "2020111B"
	t.Log(DateStrToTime(s))
}

func TestSplit(t *testing.T) {
	ss := strings.Split("", " ")
	t.Log(json.MarshalToString(ss))
	ss = ss[:0]
	t.Log(json.MarshalToString(ss))
	ss = strings.Split(" ", " ")
	t.Log(json.MarshalToString(ss))

	var a []int
	err := json.Unmarshal([]byte("[1, 2, 3]"), &a)
	t.Log(a, err)
}

func TestNil(t *testing.T) {
	var a map[int]int
	b, ok := a[123]
	t.Log(a, b, ok)

	var bs []byte
	err := json.Unmarshal(bs, &a)
	t.Log(bs, a, err, PurifyForUtf8("abc\x00xyz"))
	a = nil
	bs = []byte(`{"1":2}`)
	t.Log(json.Unmarshal(bs, &a))
	t.Log(a)
	t.Log(json.MarshalToString(a))
	var c struct {
		A int `json:"a"`
		B int `json:"b"`
	}
	bs = []byte(`{"a":2,"b":null}`)
	t.Log(json.Unmarshal(bs, &c))
	t.Log(json.MarshalToString(c))
}
