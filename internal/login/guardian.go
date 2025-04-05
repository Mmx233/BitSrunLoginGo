package login

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/config"
	"github.com/Mmx233/BitSrunLoginGo/internal/config/keys"
	"github.com/Mmx233/BitSrunLoginGo/internal/webhook"
	log "github.com/sirupsen/logrus"
	"time"
)

// Guardian 守护模式逻辑
func Guardian(logger log.FieldLogger, eventQueue webhook.EventQueue) {
	GuardianDuration := time.Duration(config.Settings.Guardian.Duration) * time.Second
	for {
		_ = Login(Conf{
			Logger:                      logger.WithField(keys.LogComponent, "login"),
			IsOnlineDetectLogDebugLevel: true,
			EventQueue:                  eventQueue,
		})
		time.Sleep(GuardianDuration)
	}
}
