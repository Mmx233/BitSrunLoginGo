# BitSrunLoginGo

深澜校园网登录脚本Go语言版。GO语言可以直接交叉编译出mips架构可执行程序（路由器）（主流平台更不用说了），从而免除安装环境。

代码逻辑来自 https://github.com/coffeehat/BIT-srun-login-script

首次运行将生成Config.json文件

Config.json说明：

```json5
{
 "from": {
  "domain": "www.msftconnecttest.com", //登录地址ip或域名
  "username": "", //账号
  "password": "" //密码
 },
 "meta": { //登录参数
  "n": "200",
  "v_type": "1",
  "acid": "5",
  "enc": "srun_bx1"
 },
 "settings": {
  "quit_if_net_ok": false, //登陆前是否检查网络
  "demo_mode": false //测试模式，报错更详细，且生成运行日志与错误日志
 }
}
```