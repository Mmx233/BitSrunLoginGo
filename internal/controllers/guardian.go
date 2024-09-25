package controllers

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/config"
	"github.com/Mmx233/BitSrunLoginGo/internal/config/keys"
	log "github.com/sirupsen/logrus"
	"time"
)

// Guardian 守护模式逻辑
func Guardian(logger log.FieldLogger) {
	GuardianDuration := time.Duration(config.Settings.Guardian.Duration) * time.Second
	for {
		_ = Login(LoginConf{
			Logger:                      logger.WithField(keys.LogComponent, "login"),
			IsOnlineDetectLogDebugLevel: true,
		})
		time.Sleep(GuardianDuration)
	}
}
