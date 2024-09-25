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
		Logger:      a.Logger,
		api:         a,
		redirectReg: redirectReg,
	}
}

type Detector struct {
	Logger log.FieldLogger

	api *Api

	redirectReg *regexp.Regexp

	// 登录页 html data
	pageUrl string
	page    []byte
}

func (d *Detector) _GET(client *http.Client, addr string) (*http.Response, error) {
	d.Logger.Debugln("HTTP GET", addr)
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range d.api.CustomHeader {
		req.Header.Set(k, fmt.Sprint(v))
	}
	return client.Do(req)
}

func (d *Detector) _DirectGET(addr string) (*http.Response, error) {
	return d._GET(d.api.Client, addr)
}

func (d *Detector) _NoDirectGET(addr string) (*http.Response, error) {
	return d._GET(d.api.NoDirect, addr)
}

func (d *Detector) _JoinRedirectLocation(addr *url.URL, loc string) (*url.URL, error) {
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

func (d *Detector) _FollowRedirect(addr *url.URL, conf _FollowRedirectConfig) (*http.Response, []byte, error) {
	addrCopy := *addr
	addr = &addrCopy

	var body []byte
	var res *http.Response
	for {
		var err error
		res, err = d._NoDirectGET(addr.String())
		if err != nil {
			return nil, nil, err
		}

		if res.StatusCode < 300 {
			body, err = io.ReadAll(res.Body)
			_ = res.Body.Close()
			if err != nil {
				return nil, nil, err
			}
			locMatch := d.redirectReg.FindSubmatch(body)
			if len(locMatch) >= 2 {
				for i := 1; i < len(locMatch); i++ {
					locBytes := locMatch[i]
					if len(locBytes) != 0 {
						addr, err = d._JoinRedirectLocation(addr, unsafe.String(unsafe.SliceData(locBytes), len(locBytes)))
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
				addr, err = d._JoinRedirectLocation(addr, res.Header.Get("location"))
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

func (d *Detector) _SearchAcid(query url.Values) (string, bool) {
	addr := query.Get(`ac_id`)
	return addr, addr != ""
}

// 用于直接获取登录页数据
func (d *Detector) _RequestPageBytes() ([]byte, error) {
	if d.pageUrl != "" {
		res, err := d._DirectGET(d.pageUrl)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			_, _ = io.Copy(io.Discard, res.Body)
			return nil, fmt.Errorf("server return http status: %d", res.StatusCode)
		}
		d.page, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return d.page, nil
	}

	baseUrl, err := url.Parse(d.api.BaseUrl)
	if err != nil {
		return nil, err
	}
	_, d.page, err = d._FollowRedirect(baseUrl, _FollowRedirectConfig{})
	return d.page, err
}

func (d *Detector) DetectEnc() (string, error) {
	if d.page == nil {
		_, err := d._RequestPageBytes()
		if err != nil {
			return "", err
		}
	}

	jsReg, err := regexp.Compile(`(?i)<script src="\.?(.+[./]portal[0-9]*\.js)(\?.*)?">`)
	if err != nil {
		return "", err
	}
	jsPathMatch := jsReg.FindSubmatch(d.page)
	if len(jsPathMatch) == 3 {
		jsPathBytes := jsPathMatch[1]
		jsPath := unsafe.String(unsafe.SliceData(jsPathBytes), len(jsPathBytes))
		jsUrl, err := url.Parse(d.api.BaseUrl)
		if err != nil {
			return "", err
		}
		jsUrl.Path = jsPath
		jsAddr := jsUrl.String()
		jsRes, err := d._DirectGET(jsAddr)
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
func (d *Detector) DetectAcid() (string, error) {
	if d.page == nil {
		// page 有值时说明 reality 已进行过 query match，此部分可跳过
		// 从入口地址 url query 寻找 acid
		baseUrl, err := url.Parse(d.api.BaseUrl)
		if err != nil {
			return "", err
		}

		var AcidFound = errors.New("acid found")
		var acid string
		_, d.page, err = d._FollowRedirect(baseUrl, _FollowRedirectConfig{
			onNextAddr: func(addr *url.URL) error {
				var ok bool
				acid, ok = d._SearchAcid(addr.Query())
				if ok {
					d.pageUrl = addr.String()
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
	result := reg.FindSubmatch(d.page)
	if len(result) == 2 {
		return string(result[1]), nil
	}

	return "", ErrAcidCannotFound
}

// Reality acid 可能为空字符串
func (d *Detector) Reality(addr string, getAcid bool) (acid string, online bool, err error) {
	startUrl, err := url.Parse(addr)
	if err != nil {
		return "", false, err
	}
	var AlreadyOnline = errors.New("already online")
	finalRes, pageBytes, err := d._FollowRedirect(startUrl, _FollowRedirectConfig{
		onNextAddr: func(addr *url.URL) error {
			if addr.Host == startUrl.Host {
				return AlreadyOnline
			}
			if getAcid {
				acid, _ = d._SearchAcid(addr.Query())
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
	d.page = pageBytes
	return
}

func (d *Detector) Reset() {
	d.pageUrl = ""
	d.page = nil
}
