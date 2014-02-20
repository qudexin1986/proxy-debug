## Introduction

This is a proxy for helping developers debug through http headers.

You can get some debugging information from http headers.

Now, just support linux, because it use shell character color.

## Why?

1. Sometimes, we want to see debugging information change when request the same url.
2. If you are debugging with a script that it will redirect to another script and redirect back.You want to see debugging information in all those script.
3. Your script have a curl request to a api designed by you, and you want to see some api's debugging information.

So, this proxy will help you! But it just can output simple debugging information, it still need to be improved.

## Install

First, you should have Golang development enviroment. Because the proxy was wrote in Golang.
Second, run `go build` to compile it.
Finally, run `ProxyDebug [-p port]`, And config the browser proxy.

## Use

require the file `PD.php`, and use:
~~~
require 'PD.php';
$name = 'rokety';
$foo  = 'Hi';
PD::info($name, $foo);
PD::warn($name);
PD::error($name);
//use group
PD::groupStart();
//many PD::info()  or PD::warn(), PD::error() method call
PD::groupEnd();//Don't forget to call
~~~

## Config

You can config the display color in config.ini and listening port.

## Screenshot

![screenshot](/screenshot.png)