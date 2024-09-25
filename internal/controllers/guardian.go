package controllers

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/config"
	"time"
)

// Guardian 守护模式逻辑
func Guardian() {
	logger := config.Logger

	logger.Infoln("[守护模式]")

	GuardianDuration := time.Duration(config.Settings.Guardian.Duration) * time.Second

	for {
		_ = Login()
		time.Sleep(GuardianDuration)
	}
}
