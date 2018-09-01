package command

import (
	"github.com/snippetor/bingo/utils"
	"runtime"
	"strings"
	"os/exec"
	"os"
)

func Run(appName, env string, b bool) {
	var execName = appName
	if runtime.GOOS == "windows" {
		execName += ".exe"
	}
	if !utils.IsFileExists(execName) {
		printError("Error: Not found app's executable file %s", execName)
		return
	}
	printInfo("Starting '%s'...", appName)
	app := appName
	if !strings.Contains(appName, "./") {
		appName = "./" + appName
	}
	var cmd *exec.Cmd
	if b {
		if runtime.GOOS == "windows" {
			printError("Error: -b(background) is not supported on windows.")
			return
		}
		args := []string{appName, ">" + app + ".log", "2>&1", "&"}
		if env == "" {
			args = append(args, "-n", app)
		} else {
			args = append(args, "-e", env, "-n", app)
		}
		cmd = exec.Command("nohup", args...)
	} else {
		var args []string
		if env == "" {
			args = []string{"-n", app}
		} else {
			args = []string{"-e", env, "-n", app}
		}
		cmd = exec.Command(appName, args...)
	}
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
