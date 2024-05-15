package config

type PgSqlConfig struct {
	Disable     bool   `toml:"disable"`
	DbConn      string `toml:"dbConn"`
	AutoMigrate bool   `toml:"autoMigrate"`
	MaxConn     int    `toml:"maxConn"`
	MaxIdleConn int    `toml:"maxIdleConn"`
	LogMode     bool   `toml:"logMode"`
}

type ServerConfig struct {
	Debug     bool   `toml:"debug"`
	DevMode   bool   `toml:"devMode"`
	RunLocal  bool   `toml:"-"` // 通过 -l flag 来设置
	CmdRoot   string `toml:"-"` // 测试时手动设置程序根目录
	Addr      string `toml:"addr"`
	JwtSecret string `toml:"jwtSecret"`
	TmpDir    string `toml:"tmpDir"`
}

type ServerLog struct {
	LogPath    string `toml:"logPath"`
	MaxSize    int    `toml:"maxSize"`
	Compress   bool   `toml:"compress"`
	MaxAge     int    `toml:"maxAge"`
	MaxBackups int    `toml:"maxBackups"`
	LogLevel   string `toml:"logLevel"`
	Format     string `toml:"format"`
}

type Workflow struct {
	XxxTheme WorkflowReq `toml:"xxxTheme"`
}

type WorkflowReq struct {
	OpId      []string `toml:"opId"`
	ServiceId int64    `toml:"serviceId"`
	Name      string   `toml:"name"`
	QueueId   string   `toml:"queueId"`
	GroupId   string   `toml:"GroupId"`
	Priority  string   `toml:"priority"`
}

type Mq struct {
	Disable    bool   `toml:"disable"`
	NameServer string `toml:"nameServer"`
	App        string `toml:"app"`
	// Ttl           int    `toml:"ttl"` // in seconds
	// ConcurrentNum int    `toml:"concurrentNum"`
	// QueSize       int    `toml:"queSize"`
	// QueTimeout    int    `toml:"queTimeout"` // in seconds
}

type Cron struct {
	DisableRecover       bool   `toml:"disableRecover"`
	DisableLoops         bool   `toml:"disableLoops"`
	ServiceCheckInterval int    `toml:"serviceCheckInterval"` // in seconds
	XxxTifDirRoot        string `toml:"xxxTifDirRoot"`
	ThemeDirRoot         string `toml:"themeDirRoot"`
}

type Ext struct {
	HttpTimeout       int               `toml:"httpTimeout"` // in seconds
	VectorStorage     string            `toml:"vectorStorage"`
	GridStorage       string            `toml:"gridStorage"`
	TileServiceCreate string            `toml:"tileServiceCreate"`
	TileServiceGet    string            `toml:"tileServiceGet"`
	XxxTifStorage     string            `toml:"xxxTifStorage"`
	XxxTifColorMap    map[string]string `toml:"xxxTifColorMap"`
	SeaweedAddr       string            `toml:"seaweedAddr"`
}

type TileService struct {
	ServiceSrids  []int32 `toml:"serviceSrids"`
	CacheEnd      int32   `toml:"cacheEnd"`
	RealTimeLayer int32   `toml:"realTimeLayer"`
	SkipNoInlay   bool    `toml:"skipNoInlay"`
}

type SelfConfig struct {
	Server   ServerConfig `toml:"server"`
	Log      ServerLog    `toml:"log"`
	Db       PgSqlConfig  `toml:"db"`
	Qp       PgSqlConfig  `toml:"qp"`
	Ds       PgSqlConfig  `toml:"ds"`
	Im       PgSqlConfig  `toml:"im"`
	Rd       PgSqlConfig  `toml:"rd"`
	Cron     Cron         `toml:"cron"`
	Mq       Mq           `toml:"mq"`
	Workflow Workflow     `toml:"workflow"`
	Ext      Ext          `toml:"ext"`
	Tile     TileService  `toml:"tile"`
}
