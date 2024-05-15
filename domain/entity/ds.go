package entity

type ShapeEntry struct {
	ShpFilePath    string   `json:"shpFilePath"`
	StorageType    string   `json:"storageType"`
	ServiceColumns []string `json:"serviceColumns"`
}

type VectorStorageResp struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    struct {
		VectorStorageId int64  `json:"vectorStorageId"`
		ShpFilePath     string `json:"shpFilePath"`
	}
}

type ParsedImgInfo struct {
	Name string `json:"Name"`
	// SatelliteID  string  `json:"satelliteId"`
	ProductLevel string `json:"productLevel"`
	// DateTime     string  `json:"DateTime"`
	Space int64  `json:"space"`
	Proj  string `json:"proj"`
	// ActualExtent string  `json:"actualExtent"`
	Extent     string  `json:"Extent"`
	AdCode     []int64 `json:"AdCode"`
	Resolution float64 `json:"Resolution"`
	BandCount  int64   `json:"BandCount"`
	BandOrder  string  `json:"BandOrder"`
}

type GridStorageReq struct {
	AdCode              []int64 `json:"adCode"`                        // 行政区划编码
	AdName              string  `json:"adName,omitempty"`              // 区域名称
	BandCount           int64   `json:"bandCount"`                     // 波段数
	BandOrder           string  `json:"bandOrder"`                     // 波段顺序
	BusinessTypeId      int64   `json:"businessTypeId,omitempty"`      // 业务类别
	CollectionDateEnd   string  `json:"collectionDateEnd,omitempty"`   // 成果采集时间结束
	CollectionDateStart string  `json:"collectionDateStart,omitempty"` // 成果采集时间开始
	DateTime            string  `json:"dateTime"`                      // 采集时间
	Extent              string  `json:"extent"`                        // 成果影像geom
	FileName            string  `json:"fileName"`                      // 影像名称
	Path                string  `json:"path"`                          // 成果影像全路径
	PathCollectionTime  string  `json:"pathCollectionTime,omitempty"`  // 文件路径上的采集时间
	ProductLevel        string  `json:"productLevel"`                  // 产品等级
	Proj                string  `json:"proj,omitempty"`                // 算子用于解析srid的字符串
	Resolution          float64 `json:"resolution"`                    // 分辨率
	Space               string  `json:"space"`                         // 坐标系
	StorageType         string  `json:"storageType"`                   // 入库方式，DATA_SERVICE
	Variogram           bool    `json:"variogram,omitempty"`           // 是否变化
	// DataType            string        `json:"dataType"`                      // 最终成果类型
	// SatelliteId         []string `json:"satelliteId"`
	// Sensor              []string `json:"sensor"`
	// SensorType          string   `json:"sensorType"`
}

type GridStorageResp struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    struct {
		GridStorageId    int64 `json:"id"`                        // 业务表id
		ResultGridMetaId int64 `json:"resultGridMetaImageInfoId"` // gp库中的id
	} `json:"data"`
}

type InitServiceParam struct {
	ServiceSource   string   `json:"service_source"`
	CutClass        string   `json:"cut_class"`
	DataClass       string   `json:"data_class"`
	MinZoom         int32    `json:"min_zoom"`
	MaxZoom         int32    `json:"max_zoom"`
	DataSetId       int64    `json:"data_set_id"`
	IsDataSet       bool     `json:"is_data_set"`
	Srid            int32    `json:"srid"`
	Srids           []int32  `json:"srids"`
	DataSetTileAddr string   `json:"data_set_tile_addr"`
	DataIds         []uint64 `json:"data_ids"`
	GridService     struct {
		GridSource string `json:"grid_source"`
		MaskFile   string `json:"maskfile"`
		Render     string `json:"render"`
		Rgba       string `json:"rgba"`
	} `json:"grid_service"`
	ServiceName         string `json:"service_name"`
	DataSetAssitanceEnd int32  `json:"data_set_assitance_end,omitempty"`
	Ts                  []Ts   `json:"ts"`
	ServiceId           int64  `json:"service_id,omitempty"`
	CacheEnd            int32  `json:"cache_end,omitempty"`
	VectorGzip          bool   `json:"vector_gzip,omitempty"`
}

type Ts struct {
	Srid        int32  `json:"srid"`         // 服务坐标系
	TilingAssis string `json:"tiling_assis"` // 切片方案辅助数据
}

type TileServiceResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    int64  `json:"data"`
}

type TileServiceStatus struct {
	Id            int64  `json:"id"`
	Status        string `json:"status"`
	IsDataSet     bool   `json:"is_data_set"`    //是否是数据集
	ServiceSource string `json:"service_source"` //data-service
	ServiceClass  string `json:"service_class"`  //grid -栅格 vector -矢量
	DataClass     string `json:"data_class"`
}
