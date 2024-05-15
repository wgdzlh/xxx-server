package entity

import (
	"errors"

	json "github.com/json-iterator/go"
)

const (
	CronTaskTableName       = "cron_tasks"
	BackgroundTaskTableName = "bg_tasks"
)

var (
	ErrWrongCronTaskGenre = errors.New("wrong cron task genre")
	ErrInvalidXxxTaskExt  = errors.New("invalid xxx task ext")
)

// 单流程任务
type CronTask struct {
	Id        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`                   // ID
	Name      string    `json:"name"`                                                 // 任务名称
	Genre     string    `json:"genre"`                                                // 任务大类型 enum(cron,download)
	Type      string    `json:"type"`                                                 // 任务子类型 enum(short,mid,factor,deduce)
	Params    AnyJson   `gorm:"type:jsonb" json:"params"`                             // 任务参数
	Dir       string    `json:"dir"`                                                  // 任务输出目录
	CreatedAt *JSONTime `gorm:"type:timestamptz(0);autoCreateTime" json:"created_at"` // 创建时间
	StartAt   *JSONTime `gorm:"type:timestamptz(0)" json:"start_at"`                  // 开始时间
	EndAt     *JSONTime `gorm:"type:timestamptz(0)" json:"end_at"`                    // 结束时间
	SrcId     uint64    `json:"src_id"`                                               // 源任务ID
	ExtId     int64     `json:"ext_id"`                                               // 关联的外部任务/流程ID
	Progress  float32   `gorm:"type:numeric(3,2)" json:"progress"`                    // 进度
	Ext       AnyJson   `gorm:"type:json" json:"ext"`                                 // 任务扩展信息
	Status    string    `gorm:"default:'NotStarted'" json:"status"`                   // 任务状态 enum(NotStarted,InProc,Done,Failed)
	ErrLog    string    `json:"err_log"`                                              // 任务错误日志
}

func (CronTask) TableName() string {
	return CronTaskTableName
}

type FactorThemeReq struct {
	Adcode        string `json:"adcode"`         // 行政区划编码
	District      string `json:"district"`       // 行政区划名称
	DistrictLevel string `json:"district_level"` // 行政区划级别，省级（province）或市级（city）
	Themes        []struct {
		Factor string `json:"factor"` // XXX要素, enum(降雨,气压,温度,湿度,风场)
		Title  string `json:"title"`  // 专题图标题
	} `json:"themes"` // 需要下载的要素列表
	ReportTime string `json:"report_time"` // 起报时刻（无需传入）
}

// 多流程任务
type BackgroundTask struct {
	Id         uint64            `gorm:"primaryKey;autoIncrement" json:"id"`                   // ID
	Type       string            `gorm:"default:'cron'" json:"type"`                           // 任务类型
	CreatedAt  *JSONTime         `gorm:"type:timestamptz(0);autoCreateTime" json:"created_at"` // 创建时间
	EndAt      *JSONTime         `gorm:"type:timestamptz(0)" json:"end_at"`                    // 结束时间
	StageEndAt *JSONTime         `gorm:"type:timestamptz(0)" json:"stage_end_at"`              // 阶段结束时间
	PubAt      *JSONTime         `gorm:"type:timestamptz(0)" json:"pub_at"`                    // 服务发布开始时间
	TaskDef    HighResInitParams `gorm:"type:json;serializer:json" json:"task_def"`            // 任务定义
	Procedure  []TaskStep        `gorm:"type:jsonb;serializer:json" json:"procedure"`          // 流程步骤状态
	Status     string            `gorm:"default:'NotStarted'" json:"status"`                   // 任务状态 enum(NotStarted,InProc,Staged,StageProc,Done,Failed)
	ErrLog     string            `json:"err_log"`                                              // 任务错误日志
}

func (BackgroundTask) TableName() string {
	return BackgroundTaskTableName
}

func (t *BackgroundTask) GetNextStep(isRestart bool) int {
	idx := -1
	for i := range t.Procedure {
		status := t.Procedure[i].Status
		if status == TASK_INIT ||
			isRestart && (status == TASK_IN_PROC || status == TASK_FAILED) {
			idx = i
			break
		}
	}
	return idx
}

func (t *BackgroundTask) GetInProcStep() int {
	idx := -1
	for i := range t.Procedure {
		if t.Procedure[i].Status == TASK_IN_PROC {
			idx = i
			break
		}
	}
	return idx
}

type TaskStep struct {
	Idx    int     `json:"-"`                // 子任务索引
	Name   string  `json:"name,omitempty"`   // 子任务名称
	Status string  `json:"status,omitempty"` // 子任务状态 enum(NotStarted,Skipped,InProc,Done,Failed)
	Ret    AnyJson `json:"ret,omitempty"`    // 子任务结果
	// Tid    []int64 `json:"tid,omitempty"`    // 子任务关联IDs
	// TaskName string `json:"task_name,omitempty"` // 子任务在其他系统中的名称
}

type HighResReq struct {
	Adcode string `json:"adcode"` // 区域adcode编码（可选，默认使用配置中的adcode）
	Year   string `json:"year"`   // 需要合成一张图的年份，例如："2022"
}

type HighResInitParams struct {
	HighResReq
	Title     string `json:"title"`      // 任务流程标题
	Manual    bool   `json:"manual"`     // 是否为手动触发数据准备流程
	PubManual bool   `json:"pub_manual"` // 是否为手动触发服务发布流程
}

type XxxTaskExt struct {
	ReportTime    JSONTime `json:"report_time"`    // 真实起报时间
	FirstForecast JSONTime `json:"first_forecast"` // 第一个预报时间
	TiffNum       int      `json:"tiff_num"`       // tiff数量
	AlertCodeNum  int      `json:"alert_code_num"` // 预警地区数量
	AlertIdxNum   int      `json:"alert_idx_num"`  // 预警网格数量
}

func (t *CronTask) GetReportTime() (rt JSONTime, err error) {
	if t.Genre != TASK_GENRE_CRON {
		err = ErrWrongCronTaskGenre
		return
	}
	if len(t.Ext) == 0 {
		err = ErrInvalidXxxTaskExt
		return
	}
	var mte XxxTaskExt
	if err = json.Unmarshal(t.Ext, &mte); err != nil {
		return
	}
	rt = mte.ReportTime
	return
}
