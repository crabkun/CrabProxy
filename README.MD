![crab](http://www.crab.pub/crabs.png)  
## CrabProxy是什么?
一个GO语言编写的支持端口转发和端口复用的小型服务端  


## CrabProxy有哪些功能？

* 端口转发
    *  例如转发SSH到HTTP，并开启HTTP混淆功能，就可以当作HTTP注射器的服务端
* 端口复用
    *  例如SSH数据包自动转发到SSH端口，HTTP数据包自动转发到HTTP端口
* HTTP混淆  
    *   支持连接建立之后第一个数据包为虚假HTTP数据包，主要用于防止某些干扰

	
## 编译

 1.  安装GOLANG编译器
    *   Centos安装命令：yum install golang
    *   Ubuntu安装命令：sudo apt-get install golang  
 2. 编译服务端主程序
    *   go build main.go
 3. 编译完成后，在源码同目录下可以得到可执行文件

 
## 配置

把下面的内容修改后保存到服务端主程序目录下，文件名为config.json
```json
    {
    "Port":"80",
    "TargetSSH":"127.0.0.1:22",
    "TargetHTTP":"127.0.0.1:801",
    "HTTPReturn200":"~crab~"
    }
```
上面示例配置的意思是，监听80端口，如果是SSH数据包，就转发到127.0.0.1:22。如果是HTTP数据包，就转发到127.0.0.1:801，如果此连接首个数据包包含\~crab\~关键字，就返回HTTP200状态码。

* Port字段为监听端口
* TargetSSH字段为转发目标SSH服务器地址和端口（可为空）
* TargetHTTP字段为转发目标HTTP服务器地址和端口（可为空）
* HTTPReturn200字段为当连接建立后第一个数据包出现某关键字时返回HTTP 200状态码数据包（可为空）  
* 注意，TargetSSH和TargetHTTP不能同时为空。HTTPReturn200字段一般用于HTTP混淆SSH，一般给HTTP注射器用，如果你没有使用，一定要把HTTPReturn200留空  


## 有问题反馈
在使用中有任何问题，欢迎反馈给我，可以用以下联系方式跟我交流

* 邮件(529493022#qq.com, 把#换成@)
* QQ: 529493022

