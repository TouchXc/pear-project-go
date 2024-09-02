package kafka

import (
	"context"
	"go.uber.org/zap"
	"ms_project/project-common/kk"
	"ms_project/project-task/internal/dao"
	"ms_project/project-task/internal/repo"
	"time"
)

var kw *kk.KafkaWriter

func InitKafkaWriter() func() {
	kw = kk.GetWriter("localhost:9092")
	return kw.Close
}
func NewCacheReader() *KafkaCache {
	reader := kk.GetReader([]string{"localhost:9092"}, "cache_group", "msproject_cache")
	return &KafkaCache{
		R:     reader,
		cache: dao.RC,
	}
}
func SendLog(data []byte) {
	kw.Send(kk.LogData{
		Topic: "msproject_log",
		Data:  data,
	})
}
func SendCache(data []byte) {
	kw.Send(kk.LogData{
		Topic: "msproject_cache",
		Data:  data,
	})
}

type KafkaCache struct {
	R     *kk.KafkaReader
	cache repo.Cache
}

func (c *KafkaCache) DelCache() {
	for {
		message, err := c.R.R.ReadMessage(context.Background())
		if err != nil {
			zap.L().Error("DelCache ReadMessage err", zap.Error(err))
			continue
		}
		zap.L().Info("收到缓存", zap.String("value", string(message.Value)))
		if "task" == string(message.Value) {
			//查询缓存中是否存在
			files, err := c.cache.HKeys(context.Background(), "task")
			if err != nil {
				zap.L().Error("DelCache HKeys err", zap.Error(err))
				continue
			}
			//删除缓存中的值
			time.Sleep(1 * time.Second)
			c.cache.Del(context.Background(), files)
		}
	}

}
