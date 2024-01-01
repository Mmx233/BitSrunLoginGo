package srun

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unsafe"
)

func (a *Api) NewDetector() *Detector {
	redirectReg, err := regexp.Compile(
		`<script>top\.self\.location\.href='(.*)'</script>|<meta http-equiv="refresh" content=".*?url=(.*?)"`,
	)
	if err != nil {
		panic(err)
	}

	return &Detector{
		api:         a,
		redirectReg: redirectReg,
	}
}

type Detector struct {
	api *Api

	redirectReg *regexp.Regexp

	// 登录页 html data
	pageUrl string
	page    []byte
}

func (a *Detector) _JoinRedirectLocation(addr *url.URL, loc string) (*url.URL, error) {
	if loc == "" {
		return nil, errors.New("目标跳转地址缺失")
	}
	if strings.HasPrefix(loc, "/") {
		locSplit := strings.Split(loc, "?")
		addr.Path = locSplit[0]
		if len(locSplit) > 1 {
			addr.RawQuery = locSplit[1]
		} else {
			addr.RawQuery = ""
		}
		return addr, nil
	} else {
		return url.Parse(loc)
	}
}

type _FollowRedirectConfig struct {
	// 获取到下一个请求地址时触发
	onNextAddr func(addr *url.URL) error
}

func (a *Detector) _FollowRedirect(addr *url.URL, conf _FollowRedirectConfig) (*http.Response, []byte, error) {
	addrCopy := *addr
	addr = &addrCopy

	var body []byte
	var res *http.Response
	for {
		log.Debugln("HTTP GET", addr)
		req, err := http.NewRequest("GET", addr.String(), nil)
		if err != nil {
			return nil, nil, err
		}
		for k, v := range a.api.CustomHeader {
			req.Header.Set(k, fmt.Sprint(v))
		}
		res, err = a.api.NoDirect.Do(req)
		if err != nil {
			return nil, nil, err
		}

		if res.StatusCode < 300 {
			body, err = io.ReadAll(res.Body)
			_ = res.Body.Close()
			if err != nil {
				return nil, nil, err
			}
			locMatch := a.redirectReg.FindSubmatch(body)
			if len(locMatch) >= 2 {
				for i := 1; i < len(locMatch); i++ {
					locBytes := locMatch[i]
					if len(locBytes) != 0 {
						addr, err = a._JoinRedirectLocation(addr, unsafe.String(unsafe.SliceData(locBytes), len(locBytes)))
						if err != nil {
							return nil, nil, err
						}
						break
					}
				}
			} else {
				break
			}
		} else {
			_, _ = io.Copy(io.Discard, res.Body)
			_ = res.Body.Close()

			if res.StatusCode < 400 {
				addr, err = a._JoinRedirectLocation(addr, res.Header.Get("location"))
				if err != nil {
					return nil, nil, err
				}
			} else {
				return nil, nil, fmt.Errorf("server return http status %d", res.StatusCode)
			}
		}

		if conf.onNextAddr != nil {
			if err = conf.onNextAddr(addr); err != nil {
				return nil, nil, err
			}
		}
	}
	return res, body, nil
}

func (a *Detector) _SearchAcid(query url.Values) (string, bool) {
	addr := query.Get(`ac_id`)
	return addr, addr != ""
}

// 用于直接获取登录页数据
func (a *Detector) _RequestPageBytes() ([]byte, error) {
	if a.pageUrl != "" {
		log.Debugln("HTTP GET", a.pageUrl)
		res, err := a.api.Client.Get(a.pageUrl)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			_, _ = io.Copy(io.Discard, res.Body)
			return nil, fmt.Errorf("server return http status: %d", res.StatusCode)
		}
		a.page, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return a.page, nil
	}

	baseUrl, err := url.Parse(a.api.BaseUrl)
	if err != nil {
		return nil, err
	}
	_, a.page, err = a._FollowRedirect(baseUrl, _FollowRedirectConfig{})
	return a.page, err
}

func (a *Detector) DetectEnc() (string, error) {
	if a.page == nil {
		_, err := a._RequestPageBytes()
		if err != nil {
			return "", err
		}
	}

	jsReg, err := regexp.Compile(`(?i)<script src="\.?(.+[./]portal[0-9]*\.js)(\?.*)?">`)
	if err != nil {
		return "", err
	}
	jsPathMatch := jsReg.FindSubmatch(a.page)
	if len(jsPathMatch) == 3 {
		jsPathBytes := jsPathMatch[1]
		jsPath := unsafe.String(unsafe.SliceData(jsPathBytes), len(jsPathBytes))
		fmt.Println("111", jsPath)
		jsUrl, err := url.Parse(a.api.BaseUrl)
		if err != nil {
			return "", err
		}
		jsUrl.Path = jsPath
		jsAddr := jsUrl.String()
		log.Debugln("HTTP GET", jsAddr)
		jsRes, err := a.api.Client.Get(jsAddr)
		if err != nil {
			return "", err
		}
		defer jsRes.Body.Close()
		if jsRes.StatusCode == 200 {
			jsContent, err := io.ReadAll(jsRes.Body)
			if err == nil {
				reg, err := regexp.Compile(`var enc = (.*?)[,;]`)
				if err != nil {
					return "", err
				}
				encMatch := reg.FindSubmatch(jsContent)
				if len(encMatch) == 2 {
					encBytes := encMatch[1]
					encStr := unsafe.String(unsafe.SliceData(encBytes), len(encBytes))
					encSplit := strings.Split(encStr, "+")
					for i, v := range encSplit {
						encSplit[i] = strings.Trim(strings.TrimSpace(v), "'\"")
					}
					enc := strings.Join(encSplit, "")
					return enc, nil
				}
			}
		} else {
			_, _ = io.Copy(io.Discard, jsRes.Body)
		}
	}

	return "", ErrEnvCannotFound
}

// DetectAcid err 为 nil 时 acid 一定存在
func (a *Detector) DetectAcid() (string, error) {
	if a.page == nil {
		// page 有值时说明 reality 已进行过 query match，此部分可跳过
		// 从入口地址 url query 寻找 acid
		baseUrl, err := url.Parse(a.api.BaseUrl)
		if err != nil {
			return "", err
		}

		var AcidFound = errors.New("acid found")
		var acid string
		_, a.page, err = a._FollowRedirect(baseUrl, _FollowRedirectConfig{
			onNextAddr: func(addr *url.URL) error {
				var ok bool
				acid, ok = a._SearchAcid(addr.Query())
				if ok {
					a.pageUrl = addr.String()
					return AcidFound
				}
				return nil
			},
		})
		if err != nil {
			if errors.Is(err, AcidFound) {
				return acid, nil
			}
			return "", err
		}
	}

	// 从 html 寻找 acid
	var reg *regexp.Regexp
	reg, err := regexp.Compile(`"ac_id".*?value="(.+)"`)
	if err != nil {
		return "", err
	}
	result := reg.FindSubmatch(a.page)
	if len(result) == 2 {
		return string(result[1]), nil
	}

	return "", ErrAcidCannotFound
}

// Reality acid 可能为空字符串
func (a *Detector) Reality(addr string, getAcid bool) (acid string, online bool, err error) {
	startUrl, err := url.Parse(addr)
	if err != nil {
		return "", false, err
	}
	var AlreadyOnline = errors.New("already online")
	finalRes, pageBytes, err := a._FollowRedirect(startUrl, _FollowRedirectConfig{
		onNextAddr: func(addr *url.URL) error {
			if addr.Host == startUrl.Host {
				return AlreadyOnline
			}
			if getAcid {
				acid, _ = a._SearchAcid(addr.Query())
			}
			return nil
		},
	})
	if err != nil {
		if errors.Is(err, AlreadyOnline) {
			online = true
			err = nil
			return
		}
		return
	}
	online = finalRes.Request.URL.Host == startUrl.Host
	a.page = pageBytes
	return
}

func (a *Detector) Reset() {
	a.pageUrl = ""
	a.page = nil
}
