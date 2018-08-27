# Bingo 重构中...

 * v0.1 重大重构
 
 * [rpcx](https://github.com/smallnest/rpcx)
 * [KCP](https://github.com/xtaci/kcp-go)
 * [etcd](https://github.com/coreos/etcd)
 
 #[完成]
 * RPC，Service消息和webApi将使用MVC模式
 * 提供中间件扩展
 * 提供模块扩展，并内置数据库模型和数据库操作模块
 
 #[未完成]
 * 增加bingo工具，提供配置文件编辑、app管理命令、调试和发布相关命令等
 * 修改配置文件格式为yaml
 * 将app转换为独立编译运行，目前是统一编译为一个可执行文件。
 * 增加监控app并提供后台管理

