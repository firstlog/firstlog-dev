## 开发中 以实现为主 未来会重构

## 配置文件
  Input:
    Tasks:
      - Task:                 # 可写多个收集任务
        Recursive: true       # 是否收集子目录 false（关闭） true（开启）
        Directory: "/data"    # 收集的目录名
        Ignore: "^\\."        # 忽略的文件名（正则匹配）
        Match: ".*."          # 收集的文件名（正则匹配）

  #  开发中    ExcludeLines: ['^DBG']          # 忽略以 ^DBG 开头的行
  #  开发中    IncludeLines: ['^ERR', '^WARN'] # 收集以 '^ERR', '^WARN' 开头的行

      - Task:
        Recursive: false
        Directory: "/data-2"
        Ignore: "^\\."
        Match: ".*."

  Output:
    Elasticsearch:
      Index: "firstlog"                  # 索引名称,启动时会检查当前索引是否存在,如果不存在的话会自动创建
      Hosts: ["http://localhost:9200"]   # elasticSearch 连接地址
      Version: "7"                       # elasticSearch 版本号,支持6.xx和7.xx, Version: "7" 则是7.xx版本
      Shards: "1"                        # 索引的分片数
      Replicas: "1"                      # 索引的副本数

      # 正则匹配收集
      Detail:
        Enable: false   # 是否开启 false（关闭） true（开启）
        Regex: '([^\s]+)\s([^\s]+)\s([^\s]+)\s\[([^\[]+)\]\s"([^"]+)"\s([^\s]+)\s([^\s]+)\s"([^"]+)"\s"([^"]+)"\s"([^"]+)"\s([^"]+)'
        Template: '{"RemoteAddr":"$1","RemoteUser":"$2","TimeLocal":"$4","Request":"$5","Status":"$6","BodyBytesSent":"$7","HttpReferer":"$8","HttpXForwardedFor":"$9","RequestTime":"$11"}'


## FirstLog Collection Tool

<img align="right" width="200px" src="https://raw.githubusercontent.com/zdwork/logo/master/firstlog.png">

FirstLog is a very lightweight tool for collecting logs. It is developed in GO (Golang) language to support multiple outputs, high performance and easy to use.

## Installation
