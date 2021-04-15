package cmd

import (
	"context"
	"fmt"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/app"
	memorystorage "github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/sirupsen/logrus"
)

type LoggerConf struct {
	Level string `json:"level"`
	Path  string `json:"path"`
}

type StorageConf struct {
	Remote   bool   `json:"remote"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Ssl      string `json:"ssl"`
}

type RabbitConf struct {
	Ttl int64  `json:"ttl"`
	Dsn string `json:"dsn"`
}

func GetStorage(ctx context.Context, log *logrus.Logger, conf StorageConf) (app.Storage, error) {
	var (
		storage app.Storage
		err     error
	)
	if conf.Remote {
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			conf.Host,
			conf.Port,
			"calendar",
			"calendar",
			conf.Database,
			conf.Ssl)
		storage, err = sqlstorage.New(ctx, log, dsn)
		if err != nil {
			return nil, err
		}
	} else {
		storage = memorystorage.New(log)
	}
	return storage, nil
}
