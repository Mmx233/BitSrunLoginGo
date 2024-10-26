package login

import (
	"context"
	"errors"
	"github.com/Mmx233/BackoffCli/backoff"
	"github.com/Mmx233/BitSrunLoginGo/internal/config"
	"github.com/Mmx233/BitSrunLoginGo/internal/config/flags"
	"github.com/Mmx233/BitSrunLoginGo/internal/config/keys"
	"github.com/Mmx233/BitSrunLoginGo/internal/dns"
	"github.com/Mmx233/BitSrunLoginGo/internal/http_client"
	"github.com/Mmx233/BitSrunLoginGo/pkg/srun"
	"github.com/Mmx233/BitSrunLoginGo/tools"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

var ipLast string
var debugTip sync.Once

type Conf struct {
	Logger                      log.FieldLogger
	IsOnlineDetectLogDebugLevel bool
}

func Login(conf Conf) error {
	logger := conf.Logger
	if config.Settings.Basic.Interfaces == "" { //单网卡
		err := Single(SingleConf{
			Conf: conf,
			Eth:  nil,
		})
		if err != nil {
			logger.Errorln("登录出错: ", err)
			debugTip.Do(func() {
				if !config.Settings.Log.DebugLevel {
					logger.Infoln("开启调试日志 (debug_level) 获取详细信息")
				}
			})
		}
		return err
	} else { //多网卡
		return Interfaces(conf)
	}
}

func ddns(logger log.FieldLogger, ip string, httpClient *http.Client) error {
	return dns.Run(&dns.Config{
		Logger:   logger.WithField(keys.LogLoginModule, "ddns"),
		Provider: config.Settings.DDNS.Provider,
		IP:       ip,
		Domain:   config.Settings.DDNS.Domain,
		TTL:      config.Settings.DDNS.TTL,
		Conf:     config.Settings.DDNS.Config,
		Http:     httpClient,
	})
}

func Interfaces(conf Conf) error {
	logger := conf.Logger
	interfaces, err := tools.GetInterfaceAddr(logger, config.Settings.Basic.Interfaces)
	if err != nil {
		return err
	}
	var interval = time.Duration(config.Settings.Basic.InterfacesInterval) * time.Second
	var errCount int
	for i, eth := range interfaces {
		logger.Infoln("使用网卡: ", eth.Name)
		if err := Single(SingleConf{
			Conf: conf,
			Eth:  &eth,
		}); err != nil {
			logger.Errorf("网卡 %s 登录出错: %v", eth.Name, err)
			errCount++
		}
		if i != len(interfaces)-1 {
			time.Sleep(interval)
		}
	}
	if errCount > 0 {
		return errors.New("multi interface login not completely succeed")
	}
	return nil
}

type SingleConf struct {
	Conf
	Eth *tools.Eth
}

func Single(conf SingleConf) error {
	if config.Settings.Backoff.Enable {
		return backoff.NewInstance(func(ctx context.Context) error {
			return doLogin(conf)
		}, config.BackoffConfig).Run(context.TODO())
	} else {
		return doLogin(conf)
	}
}

func doLogin(conf SingleConf) error {
	logger := conf.Logger

	// 登录配置初始化
	httpClient := http_client.ClientSelect(conf.Eth)
	srunClient := srun.New(&srun.Conf{
		Logger: logger,
		Https:  config.Settings.Basic.Https,
		LoginInfo: srun.LoginInfo{
			Form: *config.Form,
			Meta: *config.Meta,
		},
		Client:       httpClient,
		CustomHeader: config.Settings.CustomHeader,
	})

	srunDetector := srunClient.Api.NewDetector()

	// Reality 与 Acid
	var acidOnReality bool
	if config.Settings.Reality.Enable {
		logger := logger.WithField(keys.LogLoginModule, "reality")

		logger.Debugln("开始 Reality 流程")
		acid, _, err := srunDetector.WithLogger(logger).Reality(config.Settings.Reality.Addr, flags.AutoAcid)
		if err != nil {
			logger.Warnln("Reality 请求异常:", err)
		} else if flags.AutoAcid && acid != "" {
			acidOnReality = true
			logger.Debugf("使用嗅探 acid: %s", acid)
			srunClient.LoginInfo.Meta.Acid = acid
		}
	}
	if !acidOnReality && flags.AutoAcid {
		logger := logger.WithField(keys.LogLoginModule, "acid")

		logger.Debugln("开始嗅探")
		acid, err := srunDetector.WithLogger(logger).DetectAcid()
		if err != nil {
			if errors.Is(err, srun.ErrAcidCannotFound) {
				logger.Warnln("找不到 acid，使用配置 acid")
			} else {
				logger.Warnf("嗅探失败，使用配置 acid: %v", err)
			}
		} else {
			logger.Debugf("使用嗅探 acid: %s", acid)
			srunClient.LoginInfo.Meta.Acid = acid
		}
	}

	if flags.AutoEnc {
		logger := logger.WithField(keys.LogLoginModule, "enc")

		logger.Debugln("开始嗅探")
		enc, err := srunDetector.WithLogger(logger).DetectEnc()
		if err != nil {
			if errors.Is(err, srun.ErrEnvCannotFound) {
				logger.Warnln("找不到 enc，使用配置 enc")
			} else {
				logger.Warnf("嗅探失败，使用配置 enc: %v", err)
			}
		} else {
			logger.Debugf("使用嗅探 enc: %s", enc)
			srunClient.LoginInfo.Meta.Enc = enc
		}
	}

	// 选择输出函数
	var _Println func(args ...interface{})
	if conf.IsOnlineDetectLogDebugLevel {
		_Println = logger.Debugln
	} else {
		_Println = logger.Infoln
	}

	_Println("正在获取登录状态")

	var clientIp, loginIp string

	isClientIpRequired := !config.Meta.DoubleStack || config.Settings.DDNS.Enable
	online, ip, err := srunClient.LoginStatus()
	if err != nil {
		if isClientIpRequired && online != nil {
			logger.Debugln("响应体缺失客户端 ip，尝试从页面匹配")
			clientIp, err = srunDetector.DetectIp()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		clientIp = *ip
	}

	if config.Meta.DoubleStack {
		logger.Debugln("使用双栈网络时认证 ip 为空")
	} else {
		loginIp = clientIp
		logger.Debugln("认证客户端 ip: ", loginIp)
	}

	// 登录执行

	if *online {
		_Println("已登录~")

		if config.Settings.DDNS.Enable && config.Settings.Guardian.Enable && ipLast != clientIp {
			if ddns(logger, clientIp, httpClient) == nil {
				ipLast = clientIp
			}
		}

		return nil
	} else {
		logger.Infoln("检测到用户未登录，开始尝试登录...")

		if err = srunClient.DoLogin(loginIp); err != nil {
			return err
		}

		logger.Infoln("登录成功~")

		if config.Settings.DDNS.Enable {
			_ = ddns(logger, clientIp, httpClient)
		}
	}

	return nil
}
