package command

import (
	"bytes"
	"runtime"
	"os/exec"
	"os"
	"github.com/snippetor/bingo/config"
)

func Build(appName, env string) {
	printInfo("Start building ...")

	bingoConfig, _ := getBingoConfig(env)
	if appName == "*" {
		for _, app := range bingoConfig.Apps {
			buildApp(app)
		}
	} else {
		buildApp(bingoConfig.FindApp(appName))
	}
}

func buildApp(config *config.AppConfig) {
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
		args := []string{"build"}
		args = append(args, "-o", appName)
		cmd := exec.Command(cmdName, args...)
		cmd.Dir = config.Package
		cmd.Env = os.Environ()
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
