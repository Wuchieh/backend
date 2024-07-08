package main

import (
	"errors"
	"github.com/goccy/go-json"
	"os"
	"string_backend_0001/internal/conf"
	"string_backend_0001/internal/database"
	"string_backend_0001/internal/logger"
)

// 初始化設定
func initConfig() error {
	file, err := os.ReadFile(configFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			conf.Conf = conf.GetDefaultConfig()
			if marshal, err := json.MarshalIndent(conf.Conf, "", "  "); err != nil {
				return err
			} else {
				err = os.WriteFile(configFile, marshal, os.ModePerm)
				if err != nil {
					return err
				}
				return nil
			}
		} else {
			return err
		}
	}

	err = json.Unmarshal(file, &conf.Conf)
	if err != nil {
		return err
	}

	return nil
}

// Init 初始化
func Init() error {
	err := initConfig()
	if err != nil {
		return err
	}

	err = logger.Init()
	if err != nil {
		return err
	}

	err = database.Init()
	if err != nil {
		return err
	}
	return nil
}
