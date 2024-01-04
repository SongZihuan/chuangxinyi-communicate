# 创信易社区系统
## 介绍
创信易社区系统，是一个开源的BBS社区。提供发帖、回复、点赞、关注等功能。
此系统接入创信易身份验证系统，借助创新易系统完成身份验证。

## 如何启动
推荐docker启动。
将配置文件夹（包含config.yaml）的文件夹挂载到容器的"/usr/local/share/backend/etc"，然后挂在"/tmp/data"为可读地址，然后启动即可。
端口默认为9527。

## 推荐配置
因为配置特殊原因，并未设置统一的配置文件结构体。配置文件通过`viper.Get`在代码各处分散的被读取。

```yaml
# config.yaml
mode: debug  # gin的模式 debug或production

readableName: 创信易社区  # 名称
serviceName: 创信易社区-本地  # 日志名称

base:
  port: 9527  # 端口号
  url: http://localhost  # 服务访问的Url
  static_path: /tmp/data  # 静态文件的地址 腰围docker挂载该地址为可读

admin:
  phone: # 根管理员手机号

database:
  driver: mysql
  mysql:
    charset: utf8mb4
    host: # {mysql地址}:{mysql端口号}
    name: # 数据库名称
    user: # 用户名
    password: # 密码
    pool:
      max: 20
      min: 5
    ssl: false

redis:
  addr: # {redis地址}:{redis端口号}
  userName: default  # redis用户名
  password: # redis 密码
  db: 0 # redis数据库

jwt:
  identity_key: identity
  key: oovooYeNg1Oomisah1  # jwt的key

cors:  # 跨域设置
  allow_credentials: true
  allow_headers:
    - authorization
  allow_methods:
    - GET
    - POST
    - PUT
    - DELETE
  allow_origins:
    - http://localhost:8083
    - http://community.test.wuntsong.com
    - http://community.localhost:8083
  enable: true
  max_age: 7200

auth:
  domainUID: 61867b7a-8e38-4d00-8518-4cfc73117cdc  # 站点UID
  authDomainUID: 00000000-0000-0000-0000-000000000000  # 用户中心UID
  signPriKey: # 签名密钥（base64）
  getPubKey: # 签名公钥（base64）

  # 用户中心回调地址
  loginCaller: http://localhost:3351/api/v1/website/inspector/oauth2
  phoneCheckCaller: http://localhost:3351/api/v1/website/inspector/phone
  emailCheckCaller: http://localhost:3351/api/v1/website/inspector/email
  secondFACheckCaller: http://localhost:3351/api/v1/website/inspector/secondfa

  auditCaller: http://localhost:3351/api/v1/website/msg/audit
  workOrderCaller: http://localhost:3351/api/v1/website/msg/order
  msgCaller: http://localhost:3351/api/v1/website/msg/msg
  
  header: http://localhost:3351/api/v1/public/header/user

  websocket: ws://localhost:3351/api/v1/website/ws
  loginTokenExpire: 172800  # 用户中心的登录token有效时间

  xRunMode: release  # 运行模式

wxrobot:
  log: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=795650c9-9524-44b8-8249-60d34336244f  # 错误日志提示的企业微信webhook

aliyun:
  # 阿里云配置信息
  accessKeyId: 
  accessKeySecret: 

  # bucket信息
  bucket: # bycket名字
  endpoint: oss-cn-guangzhou.aliyuncs.com

  # IP定位API配置信息
  appKey: 
  appCode: 
  appSecret: 
  expireSecond: 86400
```

## 代码阅读指南
### 入口函数
位于`src/cmd`中，包括：

1. web是服务后端。

### 后端Api
路由位于`src/router/router.go`中，逻辑位于`src/service`中。

## 关于疑问
关于任何代码的疑问都可以向我们发送issue，我们会尽量回复。
