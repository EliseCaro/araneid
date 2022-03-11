# 天河站群系统（蜘蛛池系统）

> 一个基于Glang语言开发的站群系统（蜘蛛池系统）

### 特别说明
``` bash
  1：该项目不在更新维护，项目移步至[https://github.com/EliseCaro/eibk_client](天河站群商业版)
```

### 安装系统

``` bash
本系统采用beego框架，在开始前请确保你已经安装就绪；

拉取代码：git clone https://github.com/EliseCaro/araneid.git

进入文件夹：cd araneid

安装依赖：远行之前请安装必要的依赖库，go的依赖相对简单，哪里报错，安装那个依赖

远行：bee run

其他相关：请遵循beego文档操作
```

### 数据库配置
``` bash
基础配置在/conf/app.conf文件内，包含数据库，redis,session,百度推送等；

附加配置在数据表araneid_admin_config内；

进入后台为限制域名进入，配置项在araneid_admin_config表内admin_domain字段值；
```

### 功能版块简述
``` bash

Glang开发；实时服务器性能监控；多功能采集器；实时消息通知，附件云盘管理；

系统权限管理，拔插式配置项，第三方系统对接插件；

蜘蛛动态监控，站点数据监控，自动推送监控，站群管理，蜘蛛模型；

索引池，文档库，匹配库，关键词，栏目库，站群模板库（具备高质量多套模板）

伪原创处理（自然语言处理）；自定义监听蜘蛛，自定义监听沙盒多功能转换工具等很多功能；

```
### 简单演示图片

![](http://araneid-demo.test.upcdn.net/demo01.png)
![](http://araneid-demo.test.upcdn.net/demo02.png)

### 视频演示功能

[播放视频](http://araneid-demo.test.upcdn.net/demo.mp4)

### 如有疑问请联系本人

``` bash

# QQ交流群:682378784;SEO研究群
# 本人QQ:1368213727 (一根小腿毛)

```
