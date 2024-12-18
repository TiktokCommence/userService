package main

import (
	clog "github.com/TiktokCommence/component/log"
	"github.com/TiktokCommence/component/log/config"
	"github.com/TiktokCommence/userService/internal/conf"
	kzap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
)

func NewLogger(cf *conf.LogConf) log.Logger {
	var opts = []config.Option{
		config.WithLogFormat(config.JsonFormat),
		config.WithLogLevel(config.InfoLevel),
		config.WithLogStdout(cf.Stdout),
	}
	if cf.EnableFile {
		opts = append(opts, config.WithFileConfig(&config.FileConfig{
			LogPath:           cf.File.Path,
			LogFileName:       cf.File.Name,
			LogFileMaxSize:    int(cf.File.MaxSize),
			LogFileMaxBackups: int(cf.File.MaxBackups),
			LogMaxAge:         int(cf.File.MaxAge),
			LogCompress:       cf.File.Compress,
		}))
	}
	if cf.EnableKafka {
		opts = append(opts, config.WithKafkaConfig(&config.KafkaConfig{
			BrokersAddr: cf.Kafka.Addr,
			TopicName:   cf.Kafka.Topic,
		}))
	}

	c := config.NewConfig(opts...)
	zapLogger, err := clog.NewZapLogger(c)
	if err != nil {
		panic(err)
	}
	return kzap.NewLogger(zapLogger)
}
