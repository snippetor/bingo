package command

import (
	"os"
	"github.com/snippetor/bingo/utils"
	"github.com/snippetor/bingo/errors"
	"fmt"
)

var appTemplate = `
apps.%s = {
	package: "%s",
	domain: "127.0.0.1",
    service: {},
    rpcPort: 0,
    rpcTo: [],  
    logs: {
		default: {
            level: LevelInfo,
            outputType: OutputConsole | OutputFile,
            outputDir: ".",
            rollingType: RollingDaily
		}
	},
    db: {},
    config: {}
};
`

func New(appPackage, appName, env string) {
	var name string
	if env == "" {
		name = ".bingo.js"
	} else {
		name = ".bingo." + env + ".js"
	}
	if !utils.IsFileExists(name) {
		Init(env)
	}
	f, err := os.OpenFile(name, os.O_APPEND, 0666)
	defer f.Close()
	errors.Check(err)
	f.WriteString(fmt.Sprintf(appTemplate, appName, appPackage))
	printInfo("Bingo add %s config successfully.", appName)
}
