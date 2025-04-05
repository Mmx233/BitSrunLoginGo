package main

import (
	"github.com/Mmx233/BitSrunLoginGo/internal/config"
	"github.com/Mmx233/BitSrunLoginGo/internal/config/keys"
	"github.com/Mmx233/BitSrunLoginGo/internal/http_client"
	"github.com/Mmx233/BitSrunLoginGo/internal/login"
	"github.com/Mmx233/BitSrunLoginGo/internal/webhook"
	"time"
)

func main() {
	logger := config.Logger
	if config.Settings.Basic.Interfaces != "" {
		logger.Infoln("[多网卡模式]")
	}

	var _webhook webhook.Webhook
	if config.Settings.Webhook.Enable {
		_webhook = webhook.PostWebhook{
			Url:     config.Settings.Webhook.Url,
			Timeout: time.Duration(config.Settings.Webhook.Timeout) * time.Second,
			Client:  http_client.DefaultClient,
			Logger:  logger.WithField(keys.LogComponent, "webhook"),
		}
	} else {
		_webhook = webhook.NopWebhook{}
	}
	eventQueue := webhook.NewEventQueue(logger.WithField(keys.LogComponent, "eventQueue"), _webhook)

	if config.Settings.Guardian.Enable {
		//进入守护模式
		login.Guardian(logger.WithField(keys.LogComponent, "guard"), eventQueue)
	} else {
		//执行单次流程
		eventQueue.AddEvent(webhook.NewDataEvent(webhook.ProcessBegin, "", nil))
		_ = login.Login(login.Conf{
			Logger:                      logger.WithField(keys.LogComponent, "login"),
			IsOnlineDetectLogDebugLevel: false,
			EventQueue:                  eventQueue,
		})
		eventQueue.AddEvent(webhook.NewDataEvent(webhook.ProcessFinish, "", nil))
	}
}
