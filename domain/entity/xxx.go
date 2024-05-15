package entity

const (
	XxxDataTableName = "xxx_datas"

	DAY_TF   = "20060102"
	MID_TF   = "2006010215"
	SHORT_TF = "200601021504"
	THEME_TF = "2006-01-02 15:04"
)

type XxxData struct {
	Id       uint64 `gorm:"primaryKey;autoIncrement;comment:ID" json:"id"` // ID
	District string `gorm:"uniqueIndex" json:"district"`                   // 地区
	Path     string `json:"path"`                                          // 路径
}

func (XxxData) TableName() string {
	return XxxDataTableName
}

type XxxThemeInput struct {
	ForecastTime string   `json:"forecast_time"`
	TiffDir      string   `json:"tiff_dir"`
	Adcode       string   `json:"adcord"`
	FcType       []string `json:"fc_type"`
	Title        []string `json:"title"`
	FigureDir    string   `json:"figure_dir"`
}
