package utils

import (
	"errors"
	"strconv"
	"strings"
)

const (
	degreeSep = `°`
	minSep    = `'`
	secSep    = `"`
)

var (
	ErrMalformedDegree = errors.New("malformed degree")
)

func NorthEastToLonLat(north, east string) (lonLat [2]float64, err error) {
	lonLat[1], err = DegreeSecToFloat(north)
	if err != nil {
		return
	}
	lonLat[0], err = DegreeSecToFloat(east)
	return
}

func DegreeSecToFloat(dms string) (ret float64, err error) {
	ds := strings.SplitN(dms, degreeSep, 2)
	if len(ds) != 2 {
		err = ErrMalformedDegree
		return
	}
	ret, err = strconv.ParseFloat(strings.TrimSpace(ds[0]), 64)
	if err != nil {
		return
	}
	ds = strings.SplitN(ds[1], minSep, 2)
	if len(ds) != 2 {
		return
	}
	min, err := strconv.ParseFloat(strings.TrimSpace(ds[0]), 64)
	if err != nil {
		return
	}
	ret += min / 60.
	ds = strings.SplitN(ds[1], secSep, 2)
	if len(ds) != 2 {
		return
	}
	sec, err := strconv.ParseFloat(strings.TrimSpace(ds[0]), 64)
	if err != nil {
		return
	}
	ret += sec / 3600.
	return
}
