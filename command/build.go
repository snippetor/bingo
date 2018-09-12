package command

import (
	"bytes"
	"runtime"
	"os/exec"
	"os"
	"github.com/snippetor/bingo/config"
)

func Build(appName, env, platform string) {
	printInfo("Start building ...")

	bingoConfig, _ := getBingoConfig(env)
	if appName == "*" {
		for _, app := range bingoConfig.Apps {
			buildApp(app, platform)
		}
	} else {
		buildApp(bingoConfig.FindApp(appName), platform)
	}
}

func buildApp(config *config.AppConfig, platform string) {
	if config != nil {
		cmdName := "go"
		var (
			err    error
			stderr bytes.Buffer
		)
		appName := config.Name
		if runtime.GOOS == "windows" {
			appName += ".exe"
		}
		if platform == "" {
			platform = runtime.GOOS
		}
		args := []string{"build"}
		args = append(args, "-o", appName)
		args = append(args, "-tags", "kcp")
		cmd := exec.Command(cmdName, args...)
		cmd.Dir = config.Package
		cmd.Env = append(os.Environ(), "GOOS="+platform)
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			panic("Failed to build the application: " + stderr.String())
		}
		printSuccess("Built: %s/%s Ok.", config.Package, appName)
	} else {
		printError("Built: %s/%s failed, not found config")
	}
}
