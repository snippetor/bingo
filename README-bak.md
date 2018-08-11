# Bingo
基于golang开发的高性能，高并发分布式框架。
* 高性能，使用原生和第三方高性能包作为底，比如fasthttp, gogoprotobuf等
* 高并发，基于go语言的特性，尽量保留了并发、异步处理，减少同步处理，包括RPC，service等
* 可扩展型，每个模块都有接口定义，可以轻松实现扩展
* 可灵活配置，使用配置文件可以轻松运行多个服务；日志也可配置
* 提供多种网络及编解码方式，目前网络支持http, tcp, websocket，编解码包括json和protobuf

# 高性能组件
* http服务使用[fasthttp](https://github.com/valyala/fasthttp)，具体性能数据可以参考[HTTP server performance comparison with net/http](https://github.com/valyala/fasthttp#http-server-performance-comparison-with-nethttp)
* 编解码中protobuf使用[gogoprotobuf](https://github.com/gogo/protobuf)，并使用gogofaster模式，性能参考[Golang 序列化反序列化库的性能比较](https://github.com/smallnest/gosercomp)
* RPC使用protobuf作为消息协议，性能更高，流量更小，减少各服务之间通信消耗
* 尽量使用protobuf，不要使用json，bingo中使用的是原生json包，性能相比gogoprotobuf相差接近10倍，详细可参考[Golang 序列化反序列化库的性能比较](https://github.com/smallnest/gosercomp)

# 日志配置
```
#-------------------------------------------------------------------
# level 日志输出的最小等级
# 0: Info [默认]
# 1: Debug
# 2: Warning
# 3: Error
#-------------------------------------------------------------------
# outputType 输出类型
# 1: 控制台(Console)
# 2: 文件(File)
# 3: 控制台+文件(Console+File) [默认]
#-------------------------------------------------------------------
# logFileOutputDir 文件输出路径，默认"."
#-------------------------------------------------------------------
# logFileRollingType 日志文件分割方式
# 1: RollingDaily 按天分割一个日志文件 [默认]
# 2: RollingSize 按固定大小分割日志文件
# 3: RollingDaily+RollingSize
#-------------------------------------------------------------------
# logFileName 日志文件名，默认"bingo"
#-------------------------------------------------------------------
# logFileMaxSize 单个日志最大字节数，当包含RollingSize时生效，默认500MB
# * 可用的单位有 KB,MB,GB,TB，没有
#-------------------------------------------------------------------
# logFileScanInterval 定时扫描文件间隔，检查是否达到分割条件，单位秒，默认1秒
#-------------------------------------------------------------------
# logFileNameDatePattern 日志文件名中日期的格式，默认20060102，
# * 格式符合go标准日期格式化
#-------------------------------------------------------------------
# logFileNameExt 日志文件后缀，默认.log
#-------------------------------------------------------------------

#当前运行模式
workMode = dev

[dev]
level = 0
outputType = 3
logFileOutputDir = .
logFileRollingType = 3
logFileName = dev
logFileNameDatePattern = 20060102
logFileNameExt = .log
logFileMaxSize = 1KB
logFileScanInterval = 3

[prod]
level = 2
outputType = 2
logFileOutputDir = .
logFileRollingType = 3
logFileName = prod
logFileNameDatePattern = 20060102
logFileNameExt = .log
logFileMaxSize = 1GB
logFileScanInterval = 3
```

# 服务节点配置

```
{
  "domains": [                  -- 所有物理机内网ip地址，用于RPC通信使用
    "192.168.1.128"
  ],
  "node": [                     -- 所有服务节点
    {
      "name": "master",         -- 节点名称
      "model": "master",        -- 节点模型名称，启动时将绑定golang类型
      "domain": 0,              -- domains中的索引
      "rpc-port": 9092          -- RPC监听端口，如果是RPC服务端，需要配置此项
    },
    {
      "name": "auth",
      "model": "auth",
      "domain": 0,
      "service": [              -- 对外服务，比如网关、认证服务器需要对外服务
        {
          "name":"http8080",    -- 对外服务名称
          "type": "http",       -- 网络协议类型
          "port": 8080          -- 端口
        }
      ]
    },
    {
      "name": "gate1",
      "model": "gate",
      "domain": 0,
      "service": [
        {
          "name":"tcp9090",
          "type": "tcp",
          "port": 9090,
        }
      ],
      "rpc-port": 9091,
      "rpc-to": [               -- RPC要连接的节点，RPC客户端需要配置
        "master"
      ]
    },
    {
      "name": "game1",
      "model": "game",
      "domain": 0,
      "rpc-to": [
        "master",
        "gate1"
      ]
    }
  ]
}
```

# 启动实例

```
* echo为发布的可执行文件

# 启动echo.json中的master节点
echo start echo.json -n master 

# 启动echo.json中的所有节点，单机运行
echo start echo.json

# 命令帮助
echo -h 或者 echo
```

# TODO

 * v0.1 重大重构
 ###[完成]
 * RPC，Service消息和webApi将使用MVC模式
 * 提供中间件扩展
 * 提供模块扩展，并内置数据库模型和数据库操作模块
 ###[未完成]
 * 增加bingo工具，提供配置文件编辑、app管理命令、调试和发布相关命令等
 * 修改配置文件格式为yaml
 * 将app转换为独立编译运行，目前是统一编译为一个可执行文件。
 * 增加监控app并提供后台管理

# 示例
[bingo-example](https://github.com/snippetor/bingo-example)


# 开源协议

Bingo项目采用[Apache License v2](https://github.com/snippetor/bingo/LICENSE).发布