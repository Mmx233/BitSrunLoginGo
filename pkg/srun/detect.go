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

func (a *Api) NewDetector() Detector {
	return Detector{
		api: a,
	}
}

type Detector struct {
	api *Api

	page []byte
}

func (a *Detector) _JoinRedirectLocation(addr *url.URL, loc string) (*url.URL, error) {
	if loc == "" {
		return nil, errors.New("目标跳转地址缺失")
	}
	if strings.HasPrefix(loc, "/") {
		addr.Path = strings.TrimPrefix(loc, "/")
		return addr, nil
	} else {
		return url.Parse(loc)
	}
}

type _FollowRedirectConfig struct {
	// 覆盖响应处理逻辑，设置后 onNextAddr 无效
	onResponse func(res *http.Response) (next *url.URL, err error)
	// 获取到下一个请求地址时触发
	onNextAddr func(addr *url.URL) error
}

func (a *Detector) _FollowRedirect(addr *url.URL, conf _FollowRedirectConfig) (*url.URL, error) {
	addrCopy := *addr
	addr = &addrCopy
	for {
		log.Debugln("HTTP GET", addr)
		req, err := http.NewRequest("GET", addr.String(), nil)
		if err != nil {
			return nil, err
		}
		for k, v := range a.api.CustomHeader {
			req.Header.Set(k, fmt.Sprint(v))
		}
		res, err := a.api.NoDirect.Do(req)
		if err != nil {
			return nil, err
		}
		if conf.onResponse != nil {
			var nextAddr *url.URL
			nextAddr, err = conf.onResponse(res)
			if err != nil {
				return nil, err
			} else if nextAddr == nil {
				break
			}
			addr = nextAddr
			continue
		}
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
		if res.StatusCode < 300 {
			break
		} else if res.StatusCode < 400 {
			addr, err = a._JoinRedirectLocation(addr, res.Header.Get("location"))
			if err != nil {
				return nil, err
			}
			if conf.onNextAddr != nil {
				if err = conf.onNextAddr(addr); err != nil {
					return nil, err
				}
			}
		} else {
			return nil, fmt.Errorf("server return http status %d", res.StatusCode)
		}
	}
	return addr, nil
}

func (a *Detector) _SearchAcid(query url.Values) (string, bool) {
	addr := query.Get(`ac_id`)
	return addr, addr != ""
}

func (a *Detector) DetectEnc() (string, error) {
	log.Debugln("HTTP GET", a.api.BaseUrl)
	res, err := a.api.Client.Get(a.api.BaseUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		indexHtml, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		jsReg, err := regexp.Compile(`(?i)<script src="\.?(.+[./]portal[0-9]*\.js)(\?.*)?">`)
		if err != nil {
			return "", err
		}
		jsPathMatch := jsReg.FindSubmatch(indexHtml)
		if len(jsPathMatch) == 3 {
			jsPathBytes := jsPathMatch[1]
			jsPath := unsafe.String(unsafe.SliceData(jsPathBytes), len(jsPathBytes))
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
	} else {
		_, _ = io.Copy(io.Discard, res.Body)
	}
	return "", ErrEnvCannotFound
}

// DetectAcid err 为 nil 时 acid 一定存在
func (a *Detector) DetectAcid() (string, error) {
	// 从入口地址 url query 寻找 acid
	baseUrl, err := url.Parse(a.api.BaseUrl)
	if err != nil {
		return "", err
	}

	var AcidFound = errors.New("acid found")
	var acid string
	finalAddr, err := a._FollowRedirect(baseUrl, _FollowRedirectConfig{
		onNextAddr: func(addr *url.URL) error {
			var ok bool
			acid, ok = a._SearchAcid(addr.Query())
			if ok {
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

	// 从 html 寻找 acid
	log.Debugln("HTTP GET", finalAddr.String())
	res, err := a.api.Client.Get(a.api.BaseUrl)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		var indexHtml []byte
		indexHtml, err = io.ReadAll(res.Body)
		if err == nil {
			var reg *regexp.Regexp
			reg, err = regexp.Compile(`"ac_id".*?value="(.+)"`)
			if err != nil {
				return "", err
			}
			result := reg.FindSubmatch(indexHtml)
			if len(result) == 2 {
				return string(result[1]), nil
			}
		}
	} else {
		_, _ = io.Copy(io.Discard, res.Body)
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
	var finalUrl *url.URL
	finalUrl, err = a._FollowRedirect(startUrl, _FollowRedirectConfig{
		onResponse: func(res *http.Response) (next *url.URL, err error) {
			defer res.Body.Close()
			if res.StatusCode < 300 {
				var body []byte
				body, err = io.ReadAll(res.Body)
				if err != nil {
					return
				}

				var reg *regexp.Regexp
				reg, err = regexp.Compile(`<script>top\.self\.location\.href='(.*)'</script>`)
				if err != nil {
					return
				}

				result := reg.FindSubmatch(body)
				if len(result) == 2 {
					nextBytes := result[1]
					nextAddr := unsafe.String(unsafe.SliceData(nextBytes), len(nextBytes))
					next, err = url.Parse(nextAddr)
				}
			} else if res.StatusCode < 400 {
				next, err = a._JoinRedirectLocation(res.Request.URL, res.Header.Get("location"))
			}
			if getAcid && next != nil {
				acid, _ = a._SearchAcid(next.Query())
			}
			return
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
	online = finalUrl.Host == startUrl.Host
	return
}
