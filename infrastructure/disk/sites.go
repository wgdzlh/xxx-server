package disk

// const (
// 	siteFileName   = "tide_sites.tsv"
// 	siteFileSegNum = 5
// )

// func GetTideSitesInfo() (ts []entity.TideSiteInfo) {
// 	const tag = "GetTideSitesInfo:"
// 	raw, err := os.ReadFile(siteFileName)
// 	if err != nil {
// 		log.Fatal(tag+"ReadFile failed", zap.Error(err))
// 	}
// 	var (
// 		lines  = bytes.Split(raw, []byte{'\n'})
// 		sep    = []byte{'\t'}
// 		ss     [][]byte
// 		lonLat [2]float64
// 	)
// 	ts = make([]entity.TideSiteInfo, len(lines))
// 	for i, pair := range lines {
// 		ss = bytes.Split(pair, sep)
// 		if len(ss) != siteFileSegNum {
// 			log.Fatal(tag+"misformed tide site info: "+string(pair), zap.Int("line", i+1))
// 		}
// 		if lonLat, err = utils.NorthEastToLonLat(string(ss[3]), string(ss[4])); err != nil {
// 			log.Fatal(tag+"NorthEastToLonLat failed", zap.Error(err))
// 		}
// 		ts[i] = entity.TideSiteInfo{
// 			Id:     string(ss[0]),
// 			Region: string(ss[1]),
// 			Name:   string(ss[2]),
// 			LonLat: lonLat,
// 		}
// 	}
// 	return
// }
