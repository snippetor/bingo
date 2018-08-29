package command

import (
	"github.com/snippetor/bingo/config"
)

func Run(app, env string) {
	printInfo("Bingo start running app '%s' ...", app)
	p := config.JsParser{}
	var name string
	if env == "" {
		name = "bingo.js"
	} else {
		name = "bing." + env + ".js"
	}
	bingoConfig := p.Parse(name)

	appConfig := bingoConfig.FindApp(app)
	if appConfig == nil {
		printError("Error: Not found app by name '%s' in '%s'.", app, name)
		return
	}
}
