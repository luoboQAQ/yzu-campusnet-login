# YZU CampusNet Login

扬州大学校园网登录器

## 简介

[TerraceCN/yzu-campusnet-login](https://github.com/TerraceCN/yzu-campusnet-login)项目的Go语言重构版本

就是一个简单的登录校园网的脚本，直接下载[Releases](https://github.com/luoboQAQ/yzu-campusnet-login/releases)里构建好的二进制程序即可。

截至2023/09/27开发完成时有效。

### 起因

因为我想将该脚本放到路由器上运行，但是路由器内存有限，无法安装Python。于是决定找一门可以打包成小体积的语言来重构。一开始考虑了C++，但不会C++的网络编程😇。这时想到了Go，于是花了一天时间速通Go🙏。

### 未实现功能

由于Go的JSON解析有点麻烦，`campus_net.py`中的`get_services`函数并未实现，不过问题不大。此函数只是检验校园网服务名是否正确，由于校园网服务名一般是不会变的，所以就偷个懒🕊️。

## 用法

首先下载[Releases](https://github.com/luoboQAQ/yzu-campusnet-login/releases)里构建好的二进制程序。

然后设置环境变量或者添加`.env`文件，环境变量定义如下：

|变量名|描述|默认值|
|-|-|-|
|USER_AGENT|模拟访问所用的UA|见`config.go`|
|SSO_USERNAME|统一身份认证系统的用户名|-|
|SSO_PASSWORD|统一身份认证系统的密码|-|
|CAMPUSNET_SERVICE|校园网服务名|-|
|CHECK_INTERVAL|检测是否联网的时间间隔（秒）|60|
|START_DELAY|开始联网延时（秒）|5|
|DEBUG|调试模式（忽略是否已经连网）|false|

可用的校园网服务（依旧是截至开发完成时间）：

- 学校互联网服务
- 移动互联网服务
- 联通互联网服务
- 电信互联网服务
- 校内免费服务

`START_DELAY`不建议设的很小，由于程序出错后会自动重连，延时过小可能会导致账号被冻结。

最后直接运行程序，操作成功后程序并不会自动退出，而是会每隔60秒检测一下是否联网，断网会重连。


## 免责说明

YZU Campus Login（以下简称“本脚本”）为便于作者个人生活的脚本，本脚本所用的方法均为对正常登录过的模拟，不得用于任何商业用途。

本脚本之著作权归脚本作者所有。用户可以自由选择是否使用本脚本。如果用户下载、安装、使用本脚本，即表明用户信任该脚本作者，脚本作者对因使用项目而造成的损失不承担任何责任。