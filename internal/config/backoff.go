package config

import (
	"github.com/Mmx233/BackoffCli/backoff"
	"time"
)

var BackoffConfig backoff.Conf

func initBackoff() {
	setting := Settings.Backoff
	BackoffConfig = backoff.Conf{
		Logger:           Logger,
		DisableRecovery:  true,
		InitialDuration:  time.Duration(setting.InitialDuration) * time.Second,
		MaxDuration:      time.Duration(setting.MaxDuration) * time.Second,
		MaxRetry:         setting.MaxRetries,
		ExponentFactor:   int(setting.ExponentFactor),
		InterConstFactor: time.Duration(setting.InterConstFactor) * time.Second,
		OuterConstFactor: time.Duration(setting.OuterConstFactor) * time.Second,
	}
}
