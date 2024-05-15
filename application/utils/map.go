package utils

import (
	json "github.com/json-iterator/go"
)

func ParamsToSSMap(a any) map[string]string {
	bs, err := json.Marshal(a)
	if err != nil {
		return nil
	}
	var ret = map[string]string{}
	if err = json.Unmarshal(bs, &ret); err != nil {
		return nil
	}
	return ret
}

func ParamsToSIMap(a any) map[string]any {
	bs, err := json.Marshal(a)
	if err != nil {
		return nil
	}
	var ret = map[string]any{}
	if err = json.Unmarshal(bs, &ret); err != nil {
		return nil
	}
	return ret
}

func SplitRids(rids []int64) (pRids, nRids []int64) {
	pRids = make([]int64, 0, len(rids))
	nRids = make([]int64, 0, len(rids))
	for _, id := range rids {
		if id < 0 {
			nRids = append(nRids, -id)
		} else {
			pRids = append(pRids, id)
		}
	}
	return
}

func UniqIds(ids []uint64) []uint64 {
	rets := make([]uint64, 0, len(ids))
	ok := false
	m := map[uint64]struct{}{}
	for _, id := range ids {
		if _, ok = m[id]; !ok {
			m[id] = struct{}{}
			rets = append(rets, id)
		}
	}
	return rets
}
