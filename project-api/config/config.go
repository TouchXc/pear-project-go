package config

import (
	"github.com/spf13/viper"
	"log"
	"ms_project/project-common/logs"
	"os"
)

var InConf = InitConfig()

type Config struct {
	viper *viper.Viper
	SC    *ServerConfig
	Ec    *EtcdConfig
}
type ServerConfig struct {
	Name string
	Addr string
}
type EtcdConfig struct {
	Addrs []string
}

func InitConfig() *Config {
	conf := &Config{viper: viper.New()}
	workdir, _ := os.Getwd()
	conf.viper.SetConfigName("conf")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath(workdir + "/project-api//config")
	if err := conf.viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	conf.SC = conf.ReadServerConfig()
	conf.ReadEtcdConfig()
	return conf
}
func (c *Config) ReadServerConfig() *ServerConfig {
	sc := &ServerConfig{}
	sc.Name = c.viper.GetString("server.name")
	sc.Addr = c.viper.GetString("server.addr")
	return sc
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
