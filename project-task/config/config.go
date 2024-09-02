package config

import (
	"bytes"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"log"
	"ms_project/project-common/logs"
	"os"
)

type Config struct {
	viper       *viper.Viper
	SC          *ServerConfig
	GC          *GrpcConfig
	Ec          *EtcdConfig
	MysqlConfig *MysqlConfig
}
type ServerConfig struct {
	Name string
	Addr string
}
type GrpcConfig struct {
	Addr    string
	Name    string
	Version string
	Weight  int64
}
type EtcdConfig struct {
	Addrs []string
}
type MysqlConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DbName   string
}

func InitConfig() *Config {
	conf := &Config{viper: viper.New()}
	//先从nacos读取，如果读取不到再从本地读取
	nacosClient := InitNacosClient()
	configYaml, err := nacosClient.confClient.GetConfig(vo.ConfigParam{
		DataId: "config.yaml",
		Group:  nacosClient.group,
	})
	if err != nil {
		log.Fatalln(err)
	}
	//监听配置文件变化
	err = nacosClient.confClient.ListenConfig(vo.ConfigParam{
		DataId: "config.yaml",
		Group:  nacosClient.group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
			if err = conf.viper.ReadConfig(bytes.NewBuffer([]byte(data))); err != nil {
				log.Printf("load nacos config changed err: %s /n", err.Error())
			}
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	conf.viper.SetConfigType("yaml")
	//如果Nacos文件中存在配置文件则读取
	if configYaml != "" {
		if err = conf.viper.ReadConfig(bytes.NewBuffer([]byte(configYaml))); err != nil {
			log.Fatalln(err)
		}
	} else { //不存在则从本地读取
		workdir, _ := os.Getwd()
		conf.viper.SetConfigName("conf")
		conf.viper.AddConfigPath(workdir + "/project-task//config")
		if err = conf.viper.ReadInConfig(); err != nil {
			log.Fatalln(err)
		}
	}
	conf.SC = conf.ReadServerConfig()
	conf.GC = conf.ReadGrpcConfig()
	//conf.ReadEtcdConfig()
	conf.InitMysqlConfig()
	return conf
}

//读取user服务配置

func (c *Config) ReadServerConfig() *ServerConfig {
	return &ServerConfig{
		Name: c.viper.GetString("server.name"),
		Addr: c.viper.GetString("server.addr"),
	}
}
func (c *Config) InitZapLog() {
	//从配置中读取日志配置,初始化日志
	lc := &logs.LogConfig{
		DebugFileName: c.viper.GetString("zap.debugFileName"),
		InfoFileName:  c.viper.GetString("zap.infoFileName"),
		WarnFileName:  c.viper.GetString("zap.warnFileName"),
		MaxSize:       500,
		MaxAge:        28,
		MaxBackups:    3,
	}
	if err := logs.InitLogger(lc); err != nil {
		log.Fatalln(err)
	}
}
func (c *Config) ReadRedisConfig() *redis.Options {
	return &redis.Options{
		Addr:     c.viper.GetString("redis.host") + ":" + c.viper.GetString("redis.port"),
		Password: c.viper.GetString("redis.password"),
		DB:       c.viper.GetInt("redis.db"),
	}
}
func (c *Config) ReadGrpcConfig() *GrpcConfig {
	return &GrpcConfig{
		Addr:    c.viper.GetString("grpc.addr"),
		Name:    c.viper.GetString("grpc.name"),
		Version: c.viper.GetString("grpc.version"),
		Weight:  int64(c.viper.GetInt("grpc.weight")),
	}
}
func (c *Config) ReadEtcdConfig() {
	ec := &EtcdConfig{}
	var addrs []string
	err := c.viper.UnmarshalKey("etcd.addrs", &addrs)
	if err != nil {
		log.Fatalln(err)
	}
	ec.Addrs = addrs
	c.Ec = ec
}
func (c *Config) InitMysqlConfig() {
	mc := &MysqlConfig{
		Username: c.viper.GetString("mysql.user"),
		Password: c.viper.GetString("mysql.password"),
		Host:     c.viper.GetString("mysql.host"),
		Port:     c.viper.GetString("mysql.port"),
		DbName:   c.viper.GetString("mysql.dbname"),
	}
	c.MysqlConfig = mc
}
