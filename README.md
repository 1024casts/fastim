# FastIM

FastIM, 一个基于Go语言实现的IM服务.

> 本项目不是广义上的IM即时通信，而是具体的聊天服务

IM 服务会用到网络、数据库、缓存、加密、消息队列等，如果使用人数较多，还会涉及分布式、高并发、一致性架构设计等。

## 聊天系统组成部分

- 客户端
- 接入服务：连接保持，协议解析，session维护(标识是哪个TCP连接)，消息推送
- 业务处理服务：存储处理，消息同步，未读数等
- 存储服务：账号，消息，联系人等
- 外部接口服务(APNs，厂商服务)

## IM系统特性

- 实时性：保证消息实时触达
- 可靠性：不丢消息、消息不重复
- 一致性：多用户、多终端一致性
- 安全性：数据安全传输、数据安全存储、消息内容安全