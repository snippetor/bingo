package command

import (
	"github.com/fatih/color"
	"os"
	"io"
	"github.com/snippetor/bingo/config"
)

func getBingoConfig(env string) (*config.BingoConfig, string) {
	p := config.JsParser{}
	var name string
	if env == "" {
		name = ".bingo.js"
	} else {
		name = ".bingo." + env + ".js"
	}
	return p.Parse(name), name
}

func getAppConfig(appName, env string) *config.AppConfig {
	bingoConfig, name := getBingoConfig(env)
	appConfig := bingoConfig.FindApp(appName)
	if appConfig == nil {
		printError("Error: Not found app by name %s in config file %s.", appName, name)
		return nil
	}
	return appConfig
}

func print(format string, v ...interface{}) {
	color.Magenta(format, v...)
}

func printInfo(format string, v ...interface{}) {
	color.Blue(format, v...)
}

func printSuccess(format string, v ...interface{}) {
	color.Green(format, v...)
}

func printError(format string, v ...interface{}) {
	color.Red("ERROR: "+format, v...)
}

func copyFile(source string, dest string) error {
	ln, err := os.Readlink(source)
	if err == nil {
		return os.Symlink(ln, dest)
	}
	s, err := os.Open(source)
	if err != nil {
		return err
	}

	defer s.Close()

	d, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer d.Close()

	_, err = io.Copy(d, s)
	if err != nil {
		return err
	}

	si, err := os.Stat(source)
	if err != nil {
		return err
	}
	err = os.Chmod(dest, si.Mode())

	return err
}