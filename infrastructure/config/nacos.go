package config

import (
	"log"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

const (
	nacosAddr     = "nacos.nacos.svc.cluster.local"
	nacosPort     = 8848
	nacosUsername = "bjsh"
	nacosPassword = "pwd123"

	namespaceId = "datamanage"
	dataId      = "xxx-server"
	group       = "DEFAULT_GROUP"

	PY_CONFIG_SEP = "---"
)

// 从Nacos获取配置
func GetConfigFromNacos() (content string) {
	sc := []constant.ServerConfig{
		{
			IpAddr: nacosAddr,
			Port:   nacosPort,
		},
	}

	cc := &constant.ClientConfig{
		NotLoadCacheAtStart: true,
		LogLevel:            "warn",
		NamespaceId:         namespaceId,
		Username:            nacosUsername,
		Password:            nacosPassword,
	}

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		log.Println("make nacos client error:", err.Error())
		return
	}

	content, err = client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		log.Println("nacos get config error:", err.Error())
		return ""
	}

	listenConfig(client)
	return
}

// 监听配置变更
func listenConfig(client config_client.IConfigClient) {
	err := client.ListenConfig(vo.ConfigParam{
		DataId:   dataId,
		Group:    group,
		OnChange: updateConfig,
	})
	if err != nil {
		log.Println("nacos listen config error:", err.Error())
	}
}

// 更新配置
func updateConfig(namespace, group, dataId, data string) {
	log.Println("nacos config changed, namespace:", namespace, "group:", group, "dataId:", dataId)
	subs := strings.SplitN(data, PY_CONFIG_SEP, 2)
	if len(subs) == 0 || subs[0] == "" {
		log.Println("new nacos config is empty, discard")
		return
	}
	newCfg := &SelfConfig{}
	if _, err := toml.Decode(subs[0], newCfg); err != nil {
		log.Println("fail to decode config:", err.Error())
		return
	}
	setGlobalConfig(newCfg)
}
