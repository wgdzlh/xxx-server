package cmd

import (
	slog "log"
	"os"
	"path/filepath"
	"strings"

	"xxx-server/application/app"
	log "xxx-server/application/logger"
	"xxx-server/application/utils"
	"xxx-server/infrastructure/config"
	"xxx-server/infrastructure/persistence"
	"xxx-server/infrastructure/shell"
	"xxx-server/interface/api"

	"github.com/BurntSushi/toml"
	json "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func SetupConfig() {
	var data string
	if config.C.Server.RunLocal {
		tFile := filepath.Join(config.C.Server.CmdRoot, "config.toml")
		cfg, _ := os.ReadFile(tFile)
		data = string(cfg)
	} else {
		data = config.GetConfigFromNacos()
	}
	subs := strings.SplitN(data, config.PY_CONFIG_SEP, 2)
	if len(subs) > 0 && subs[0] != "" {
		if _, err := toml.Decode(subs[0], config.C); err != nil {
			slog.Fatal("config decode failed:", err)
		}
		if len(subs) < 2 || subs[1] == "" {
			return
		}
		if err := os.WriteFile(filepath.Join(config.C.Server.CmdRoot, "config.ini"),
			utils.S2B(subs[1]), os.ModePerm); err != nil {
			slog.Fatal("failed to output python config.ini:", err)
		}
	} else {
		slog.Fatal("config not initialized")
	}
}

func Execute() {
	cobra.OnInitialize(SetupConfig)

	rootCmd := &cobra.Command{
		Use:   "xxx-server",
		Short: "Start Run xxx-server Server",
		Run:   serverRun,
	}
	rootCmd.Flags().BoolVarP(&config.C.Server.RunLocal, "local", "l", false, "test in local environment")

	if err := rootCmd.Execute(); err != nil {
		slog.Println("rootCmd execute failed:", err)
	}
}

func serverRun(cmd *cobra.Command, args []string) {
	cfgStr, _ := json.MarshalToString(config.C)
	slog.Println("config:", cfgStr)
	logCfg := config.C.Log
	log.InitLog(
		log.Path(logCfg.LogPath),
		log.Level(logCfg.LogLevel),
		log.Compress(logCfg.Compress),
		log.MaxSize(logCfg.MaxSize),
		log.MaxBackups(logCfg.MaxBackups),
		log.MaxAge(logCfg.MaxAge),
		log.Format(logCfg.Format),
	)
	config.AddCallback(func(sc *config.SelfConfig) {
		newLevel := sc.Log.LogLevel
		if newLevel != config.C.Log.LogLevel {
			log.Info("set new log level", zap.String("level", newLevel))
			log.SetLevel(newLevel)
		}
	})

	err := persistence.SetupRepositories()
	if err != nil {
		log.Fatal("initDb failed", zap.Error(err))
	}

	// err = shell.SetupRepositories()
	// if err != nil {
	// 	log.Fatal("setup shell failed", zap.Error(err))
	// }

	app.Init(persistence.R, shell.R) // 注入底层数据库依赖

	// 测试模式，不需要起API
	// if cmd == nil {
	// return
	// }

	err = api.Run()
	if err != nil {
		log.Error("apiRun error", zap.Error(err))
	}
}
