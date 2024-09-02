package config

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
)

type NacosClient struct {
	confClient config_client.IConfigClient
	group      string
}

func InitNacosClient() *NacosClient {
	bootstrapConf := InitBootstrap()
	//创建客户端配置 clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         bootstrapConf.NacosConfig.Namespace, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	//创建服务端配置  ServerConfig至少一个
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      bootstrapConf.NacosConfig.IpAddr,
			ContextPath: bootstrapConf.NacosConfig.ContextPath,
			Port:        uint64(bootstrapConf.NacosConfig.Port),
			Scheme:      bootstrapConf.NacosConfig.Scheme,
		},
	}
	// 创建动态配置 NewConfigClient
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	nc := &NacosClient{
		confClient: configClient,
		group:      bootstrapConf.NacosConfig.Group,
	}
	return nc
}
