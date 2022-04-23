# NWAFU-SrunLoginGo
[![Lisense](https://img.shields.io/github/license/Mmx233/BitSrunLoginGo)](https://github.com/Mmx233/BitSrunLoginGo/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/Mmx233/BitSrunLoginGo?color=blueviolet&include_prereleases)](https://github.com/Mmx233/BitSrunLoginGo/releases)
[![GoReport](https://goreportcard.com/badge/github.com/Mmx233/BitSrunLoginGo)](https://goreportcard.com/report/github.com/Mmx233/BitSrunLoginGo)
## 本项目仅对配置文件做了补充修改，经测试可用于西北农林科技大学深澜校园网网页认证。包括Windows、Linux(路由器padavan,OpenWrt系统同样适用)
编译详情请查看 https://github.com/Mmx233/BitSrunLoginGo 的项目，本项目只介绍使用方法

> 主要登录逻辑来自： https://github.com/coffeehat/BIT-srun-login-script

> 对Openwrt更加友好的ipk编译版： [Mmx233/BitSrunLoginGo_Openwrt](https://github.com/Mmx233/BitSrunLoginGo_Openwrt) ，该版本压缩了binary文件，节省闪存空间


## 使用方法
1. 根据设备架构下载编译后的可执行文件，解压后获得名为`autoLogin`的文件
   > [点这里跳转到原项目的可执行文件下载页面](https://github.com/Mmx233/BitSrunLoginGo/releases)
2. 在`autoLogin`同级文件下创建`Config.yaml`文件
    > 可以通过添加启动参数`--config`指定配置文件路径，默认为当前目录的`Config.yaml`。
    
    > 建议仅使用`json`和`yaml`
3. 修改`Config.yaml`的内容如下（复制后更改其中的 `学号` 和 `密码` 即可）：
    ```yaml
    form:
        domain: portal.nwafu.edu.cn
        username: "学号"
        usertype:
        password: "密码"
    meta:
        "n": "200"
        type: "1"
        acid: "1"
        enc: srun_bx1
    settings:
        basic:
            https: true
            skip_cert_verify: false
            timeout: 5
            interfaces: ""
        guardian:
            enable: false
            duration: 300
        daemon:
            enable: false
            path: .autoLogin
        debug:
            enable: false
            write_log: false
            log_path: ./
    ```
4. 运行`autoLogin`程序
   > Windows平台下可直接双击，或在当前目录打开cmd，输入```autoLogin```运行

   > Linux下，先使用```chmod a+x autoLogin```赋予可执行权限，再执行```./autoLogin```

至此应该认证完成，打开 https://portal.nwafu.edu.cn/ 应该会显示已登录页面

## 附 使用`--config`指定配置文件

```shell
./autoLogin --config=/pathA/config1.yaml
./autoLogin --config=/pathB/config2.yaml
```

首次运行将自动生成配置文件

`Config.yaml`中各字段的说明，以`Config.json`为例：

```json5
{
  "form": {
    "domain": "www.msftconnecttest.com", //登录地址ip或域名
    "username": "", //账号
    "user_type": "cmcc", //运营商类型，西农此值为空
    "password": "" //密码
  },
  "meta": { //登录参数
    "n": "200",
    "type": "1",
    "acid": "5",
    "enc": "srun_bx1"
  },
  "settings": {
    "basic": { //基础设置
      "https": false, //访问校园网API时直接使用https URL
      "skip_cert_verify": false, //是否忽略证书验证
      "interfaces": "", //网卡名称正则（注意JSON转义），如：eth0\\.[2-3]，不为空时为多网卡模式
      "timeout": 5 //网络请求超时时间（秒）
    },
    "guardian": { //守护模式
      "enable": false,
      "duration": 300, //网络检查周期（秒）
    }, 
    "daemon": { //后台模式（不建议windows使用）
      "enable": false,
      "path": ".BitSrun", //守护监听文件路径，确保只有单守护运行
    },
    "debug": {
      "enable": false, //开启debug模式，报错将更加详细
      "write_log": false, //写日志文件
      "log_path": "./" //日志文件存放路径
    }
  }
}
```

登录参数从原网页登陆时对`/srun_portal`的请求抓取，抓取时请把浏览器控制台的`preserve log`（保留日志）启用。
