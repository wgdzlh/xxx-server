package entity

type WorkflowReq struct {
	TaskId            string         `json:"-"`
	TaskType          string         `json:"-"`
	WorkflowServiceId int64          `json:"workflowServiceId"`
	Name              string         `json:"name"`
	QueueId           string         `json:"queueId"`
	GroupId           string         `json:"groupId"`
	Priority          string         `json:"priority"`
	Input             WorkflowInputs `json:"input,omitempty"`
	Param             WorkflowParams `json:"param,omitempty"`
}

type WorkflowInputs = []map[string]any
type WorkflowParams = []any

type WorkflowMqReq struct {
	*WorkflowReq
	Ext MQInfo `json:"mqInfo"`
}

type MQInfo struct {
	Topic       string   `json:"topic"`
	Tag         string   `json:"tag,omitempty"`
	CustomMsgId string   `json:"customMsgId"`
	Keys        []string `json:"keys"`
}

type WorkflowMqResp struct {
	CustomMsgId  string  `json:"customMsgId"`
	AssignmentId int     `json:"assignmentId"`
	Code         int     `json:"code"`
	Msg          string  `json:"msg"`
	Data         AnyJson `json:"data"`
}

type TaskResult struct {
	Id       int64   `json:"id"`
	Progress float32 `json:"progress"`
}

type XxxTifStoreReq struct {
	// ImageType           string       `json:"ImageType"`
	FileName            string       `json:"fileName"`
	FilePath            string       `json:"filePath"`
	StorageType         string       `json:"storageType"`
	FactorType          string       `json:"factorType"`
	ReportTime          JSONTime     `json:"reportTime"`
	DeduceTime          JSONTime     `json:"deduceTime"`
	DeduceSpan          string       `json:"deduceSpan"`
	Srid                int32        `json:"srid"`
	BandCount           int32        `json:"bandCount"`
	BandsOrder          string       `json:"bandsOrder"`
	ResourceName        string       `json:"resourceName"`
	ResourceDesc        string       `json:"resourceDesc"`
	VirtualDatabaseName string       `json:"virtualDatabaseName"`
	Resolution          float32      `json:"resolution"`
	AdCode              []int32      `json:"adCode"`
	Extent              string       `json:"extent"`
	RenderOption        RenderOption `json:"renderOption"`
	IsTileService       bool         `json:"isTileService"`
	// VirtualDatabaseID   int32        `json:"virtualDatabaseId"`
}

type RenderOption struct {
	Render            string   `json:"render"`
	Rgba              string   `json:"rgba"`
	RangeColorMap     string   `json:"rangecolormap"`
	Stretch           string   `json:"strech,omitempty"`
	CumulativeCutLow  *float64 `json:"cumulativecutlow,omitempty"`
	CumulativeCutHigh *float64 `json:"cumulativecuthigh,omitempty"`
}
