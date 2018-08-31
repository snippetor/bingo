package command

import (
	"bytes"
	"os"
	"runtime"
	"os/exec"
)

func Build(appPackage, appName string) {
	printInfo("Start rebuilding ...")
	cmdName := "go"
	var (
		err    error
		stderr bytes.Buffer
	)
	if runtime.GOOS == "windows" {
		appName += ".exe"
	}
	args := []string{"build"}
	args = append(args, "-o", appName)
	cmd := exec.Command(cmdName, args...)
	cmd.Dir = appPackage
	cmd.Env = append(os.Environ(), "GOGC=off")
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		printError("Failed to build the application: %s", stderr.String())
		return
	}
	printInfo("Built Successfully!")
}
