package app

import (
	"github.com/ZAF07/tigerlily-e-bakery-payment/internal/config"
	"github.com/ZAF07/tigerlily-e-bakery-payment/internal/db"
	"github.com/ZAF07/tigerlily-e-bakery-payment/internal/pkg/constants"
	"github.com/ZAF07/tigerlily-e-bakery-payment/internal/pkg/logger"
)

func InitApplication() {
	appConfig := config.LoadConfig()
	initDB(appConfig)
	initLogger(appConfig)
}

func initDB(appConfig *config.ApplicationConfig) {

	switch appConfig.GeneralConfig.DBType {
	case constants.POSTGRES:
		appConfig.PaymentDB = db.InitPostgresDB()
	default:
		db.NewDB()
	}
}

func initLogger(appConfig *config.ApplicationConfig) {
	loggerType := appConfig.GeneralConfig.Logger

	switch loggerType {
	case constants.Default:
		appConfig.DefaultLogger = loadDefaultLogger()
	}
	return
}

func loadDefaultLogger() *logger.Logger {
	return logger.NewLogger()
}