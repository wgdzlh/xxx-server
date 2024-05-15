package entity

const (
	OptionTableName = "options"

	SepInKey   = "@"
	SepInValue = ","

	OPT_TRUE   = "1"
	OPT_FALSE  = "0"
	OPT_FAILED = "-1"
)

type Option struct {
	Key   string `gorm:"primaryKey" json:"key"`
	Value string `json:"value"`
}

func (Option) TableName() string {
	return OptionTableName
}
