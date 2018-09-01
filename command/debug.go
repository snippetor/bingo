package command

func Debug(appName, env string) {
	appConfig := getAppConfig(appName, env)
	if appConfig != nil {
		Watch(appConfig.Package, appName, env)
	}
}
