package logger

import (
	"os"

	
	"go.uber.org/zap"
)
var ZapLogger *zap.Logger

func StartLogger() (error) {
	mode := os.Getenv("APP_ENV")
	var err error
	if mode == "" || mode == "DEV" {
		ZapLogger, err = zap.NewDevelopment()
		if err != nil {
			return err
		}
	} else {
		ZapLogger, err = zap.NewProduction()
		if err != nil {
			return err
		}
	}
	ZapLogger.Info("zap logger is running!")
	return  nil
}