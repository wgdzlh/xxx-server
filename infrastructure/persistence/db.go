package persistence

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "xxx-server/application/logger"
	"xxx-server/domain/entity"
	"xxx-server/domain/repository"
	"xxx-server/infrastructure/config"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Repositories struct {
	db       *gorm.DB // 本服务后端DB
	ds       *gorm.DB // 数仓业务DB
	im       *gorm.DB // GP影像存储DB
	OutSvr   repository.OutSvrRepo
	Option   repository.OptionRepo
	Setting  repository.SettingRepo
	CronTask repository.CronTaskRepo
	XxxData  repository.XxxDataRepo
}

var (
	R *Repositories

	validInts = regexp.MustCompile(`^[0-9,-]+$`)

	ErrDbNotSet = errors.New("dbConn not set")
)

func SetupRepositories() (err error) {
	R = &Repositories{}
	if R.db, err = getDB(config.C.Db); err != nil {
		return
	}
	if R.ds, err = getDB(config.C.Ds); err != nil {
		return
	}
	if R.im, err = getDB(config.C.Im); err != nil {
		return
	}
	if err = autoMigrate(); err != nil {
		log.Error("migrateDb failed", zap.Error(err))
		return
	}
	R.OutSvr = NewOutSvrRepo()
	R.Option = NewOptionRepo()
	R.Setting = NewSettingRepo()
	R.CronTask = NewCronTaskRepo()
	R.XxxData = NewXxxDataRepo()
	config.AddCallback(R.syncLogMode)
	log.Info("setup databases succeed")
	return
}

func getDB(cfg config.PgSqlConfig) (db *gorm.DB, err error) {
	if cfg.Disable {
		return
	}
	log.Info("database config", zap.String("dsn", cfg.DbConn))
	if cfg.DbConn == "" {
		err = ErrDbNotSet
		return
	}
	dbCfg := &gorm.Config{
		PrepareStmt: true,
	}
	if cfg.LogMode {
		dbCfg.Logger = logger.Default.LogMode(logger.Info)
	} else {
		dbCfg.Logger = logger.Default.LogMode(logger.Error)
	}
	db, err = gorm.Open(postgres.Open(cfg.DbConn), dbCfg)
	if err != nil {
		return
	}
	dbPool, err := db.DB()
	if err != nil {
		return
	}
	oc, ic := 100, 5
	if cfg.MaxConn > 0 {
		oc = cfg.MaxConn
	}
	if cfg.MaxIdleConn > 0 {
		ic = cfg.MaxIdleConn
	}
	dbPool.SetMaxOpenConns(oc)
	dbPool.SetMaxIdleConns(ic)
	dbPool.SetConnMaxLifetime(time.Hour)

	err = testDBConnect(db)
	return
}

func testDBConnect(db *gorm.DB) (err error) {
	var version string
	if err = db.Raw("SELECT version()").Scan(&version).Error; err != nil {
		return errors.New("db connect failed: " + err.Error())
	}
	log.Info("database connected", zap.String("version", version))
	return
}

func (r *Repositories) syncLogMode(sc *config.SelfConfig) {
	oc := config.C
	if sc.Db.LogMode != oc.Db.LogMode {
		setDbLogMode(r.db, sc.Db.LogMode)
	}
	if sc.Ds.LogMode != oc.Ds.LogMode {
		setDbLogMode(r.ds, sc.Ds.LogMode)
	}
	if sc.Im.LogMode != oc.Im.LogMode {
		setDbLogMode(r.im, sc.Im.LogMode)
	}
}

func setDbLogMode(db *gorm.DB, logMode bool) {
	if db == nil {
		return
	}
	if logMode {
		db.Logger = db.Logger.LogMode(logger.Info)
	} else {
		db.Logger = db.Logger.LogMode(logger.Error)
	}
}

// 同步更新所有表
func autoMigrate() (err error) {
	if config.C.Db.AutoMigrate {
		entities := []any{
			&entity.Option{},
			&entity.Setting{},
			&entity.CronTask{},
			&entity.XxxData{},
		}
		if err = R.db.AutoMigrate(entities...); err != nil {
			return
		}
		// if err = R.db.Exec(`DELETE FROM "`+entity.ResourceDirTableName+`" WHERE type = ? AND sub_type NOT IN ?`,
		// 	entity.RES_TYPE_IMG, []string{entity.SUB_TYPE_RET, entity.SUB_TYPE_DOM}).Error; err != nil {
		// 	return
		// }
	}
	return
}

func ToIntArray(ids []int64) clause.Expr {
	var ret strings.Builder
	ret.WriteString("ARRAY[")
	for i, id := range ids {
		if i > 0 {
			ret.WriteByte(',')
		}
		ret.WriteString(strconv.FormatInt(id, 10))
	}
	if len(ids) == 0 {
		ret.WriteString("]::int[]")
	} else {
		ret.WriteByte(']')
	}
	return gorm.Expr(ret.String())
}

func Int32ToIntArray(ids []int32) clause.Expr {
	var ret strings.Builder
	ret.WriteString("ARRAY[")
	for i, id := range ids {
		if i > 0 {
			ret.WriteByte(',')
		}
		ret.WriteString(strconv.FormatInt(int64(id), 10))
	}
	if len(ids) == 0 {
		ret.WriteString("]::int[]")
	} else {
		ret.WriteByte(']')
	}
	return gorm.Expr(ret.String())
}

func ToBigIntArray(ids []int64) clause.Expr {
	var ret strings.Builder
	ret.WriteString("ARRAY[")
	for i, id := range ids {
		if i > 0 {
			ret.WriteByte(',')
		}
		ret.WriteString(strconv.FormatInt(id, 10))
	}
	ret.WriteString("]::int8[]")
	return gorm.Expr(ret.String())
}

func ToStringArray(ss []string) clause.Expr {
	var ret strings.Builder
	ret.WriteString("ARRAY[")
	for i, s := range ss {
		if i > 0 {
			ret.WriteByte(',')
		}
		ret.WriteByte('\'')
		ret.WriteString(s)
		ret.WriteByte('\'')
	}
	ret.WriteByte(']')
	return gorm.Expr(ret.String())
}

func Int32SpansToOrArray(ids [][2]int32, field string) clause.Expr {
	template := field + " BETWEEN %d AND %d"
	var ret strings.Builder
	for i, id := range ids {
		if i > 0 {
			ret.WriteString(" OR ")
		}
		ret.WriteString(fmt.Sprintf(template, id[0], id[1]))
	}
	return gorm.Expr(ret.String())
}

func TransSQLError(err error) error {
	msg := err.Error()
	var aMsg string
	switch {
	case strings.HasSuffix(msg, "(SQLSTATE 23503)"):
		aMsg = "字段不符合外键约束 - " + msg
	default:
		aMsg = "未知SQL错误 - " + msg
	}
	return errors.New(aMsg)
}

// anti SQL injection
func CheckSQLIntsValue(s string) bool {
	return validInts.MatchString(s)
}
