package entity

const (
	SettingTableName = "settings"

	SECTION_ALERT = "alert"
)

type Setting struct {
	Id      uint64  `gorm:"primaryKey;autoIncrement" json:"id"`                     // ID
	Section string  `gorm:"uniqueIndex:section-with-value" json:"section"`          // 设置所属分栏，enum(预警-alert)
	Value   AnyJson `gorm:"type:jsonb;uniqueIndex:section-with-value" json:"value"` // 设置的值
}

func (Setting) TableName() string {
	return SettingTableName
}

type AlertDistConfig struct {
	Adcode   string    `json:"adcode"`                    // 行政区划编码
	District string    `json:"district"`                  // 行政区划名称
	LatLon   []float64 `json:"lat_lon"`                   // 经纬度范围，排序：经度最小值，经度最大值，纬度最小值，纬度最大值
	Subs     []int32   `json:"subs" swaggerignore:"true"` // 行政区划管辖的区县级编码
	Wkt      string    `json:"wkt,omitempty"`             // 输入范围的WKT
}
