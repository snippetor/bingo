package command

import (
	"github.com/snippetor/bingo/config"
)

func Run(appName, env string) {
	p := config.JsParser{}
	var name string
	if env == "" {
		name = ".bingo.js"
	} else {
		name = ".bingo." + env + ".js"
	}
	bingoConfig := p.Parse(name)
	appConfig := bingoConfig.FindApp(appName)
	if appConfig == nil {
		printError("Error: Not found app by name '%s' in '%s'.", appName, name)
		return
	}
	watch(appConfig.Package, appName, env)
}
