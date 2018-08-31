package command

import (
	"io/ioutil"
	"github.com/snippetor/bingo/utils"
)

var template = `/**
 * 常量定义
 */
// 日志级别
LevelInfo = 0;
LevelDebug = 1;
LevelWarning = 2;
LevelError = 3;
// 日志输出类型
OutputConsole = 1 << 0;
OutputFile = 1 << 1;
// 日志分割选项
RollingNone = 0;
RollingDaily = 1 << 0;
RollingSize = 1 << 1;
// 日志最大字节数
KB = 1 << 10;
MB = 1 << 20;
GB = 1 << 30;
TB = 1 << 40;
// 编解码
Json = "json";
Protobuf = "protobuf";
// 网络类型
Http = "http";
Kcp = "kcp";
Tcp = "tcp";
Websocket = "ws";
// 数据库
MySql = "mysql";
Mongo = "mongo";

/**
 * 全局配置
 *
 * @enableBingoLog: 是否显示框架日志
 */
config = {
    enableBingoLog: true,
};
/**
 * apps中定义所有节点配置。
 * 单个节点配置如下：
 * apps.app1 = {
 *   package: "app",
 *   etcds: [],
 *   service: {
 *      http8080: { net: "http", port: 8080, codec: "json" },
 *      kcp9090: { net: "kcp", port: 9090 },
 *      tcp9091: { net: "tcp", port: 9091 },
 *      ws9092: { net: "ws", port: 9092 },
 *   },
 *   rpcPort: 0,           
 *   rpcTo: ["app2"],      
 *   logs: {
 *       default: {
 *           level: LevelInfo,
 *           outputType: OutputConsole | OutputFile,
 *           outputDir: ".",
 *           rollingType: RollingDaily | RollingSize,
 *           fileMaxSize: 500*MB
 *       },
 *   },
 *   db: {
 *      mongo1: {
 *          type: "mongo"
 *          addr: "localhost:27017"
 *          user: "",
 *          pwd: "",
 *          db: "test",
 *      },
 *      mysql1: {
 *          type: "mysql"
 *          addr: "localhost:27017"
 *          user: "",
 *          pwd: "",
 *          db: "test",
 *          tbPrefix: "tb"
 *      }
 *   },
 *   config: {}
 *  }
 *  @app1: 为app唯一名称，可自定义，并在运行时作为进程名，日志名称也使用它
 *  @package：节点代码实现所在的工程目录
 *  @etcds: 设置需要注册的etcd服务器地址
 *  @service: 定义节点需要启用的服务，
 *      @net: 可以为http, kcp, tcp, ws
 *      @codec: 默认为protobuf, 可以为json, protobuf
 *  @rpcPort: 如果需要启用RPC Server则设置此端口；反之则设置为0或不设置
 *  @rpcTo: 如果需要RPC链接到那个app则添加到rpcTo；反之则设置为空或不设置
 *  @log: 日志配置，必须设置default日志，如果不设置日志将只输出到控制台
 *      @level: 日志打印的最小级别
 *      @outputType: 输出类型
 *      @outputDir: 输出位置
 *      @rollingType: 日志分割选项
 *      @fileMaxSize: 如果设置RollingSize，则表示日志文件大小；否则无效
 *  @db: 数据库配置，可以关联上面定义的dbs;
 *      @type: 目前支持mongo, mysql
 *      @tbPrefix: 数据库表前缀，mysql可以设置
 *  @config: 自定义配置
 */
apps = {};
`
func Init(env string) {
	var name string
	if env == "" {
		name = ".bingo.js"
	} else {
		name = ".bingo." + env + ".js"
	}
	if utils.IsFileExists(name) {
		printError("Bingo init failed, %s is exists.", name)
		return
	}
	printInfo("Bingo init config file %s... ", name)
	ioutil.WriteFile(name, []byte(template), 0666)
	printSuccess("Bingo init Done!")
}
