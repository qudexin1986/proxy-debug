## 介绍

这个代理程序主要是为了帮助开发这可以将调试信息输出到HTTP头部。

现在，只支持Linux系统，因为它使用了shell的字符颜色功能。


## 为什么要写这个代理程序

1. 有时，我们想看到对于同一个URL请求多次，输出的调试信息的变化。
2. 或者当你调试的脚本需要重定向到其他脚本，然后又重定向回来，你希望能看到这些脚本的调试输出信息。
3. 当你的脚本通过CURL请求一个自己设计的API时，你想输出这个API的一些调试信息。

那么这个代理程序将能够帮到你，但它现在只能简单的输出一些调试信息，它仍需要完善。

## 安装

首先，你需要有Golang的开发环境，因为这个代理是使用Golang写的。

然后，执行`go build`来编译它。

最后，执行`ProxyDebug [-p port]`，并设置好浏览器代理。

## Use

require the file `PD.php`, and use:
~~~
require 'PD.php';
$name = 'rokety';
PD::info($name);
PD::warn($name);
PD::error($name);
//use group
PD::groupStart();
//many PD::info()  or PD::warn(), PD::error() method call
PD::groupEnd();//Don't forget to call
~~~

## 配置

你可以配置字符的颜色和监听的端口。

## 截图

![screenshot](/screenshot.png)