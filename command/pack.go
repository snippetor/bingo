package command

import (
	"os"
	"runtime"
	"path/filepath"
	"io"
	"archive/zip"
	"time"
	"github.com/snippetor/bingo/utils"
	"fmt"
)

func Pack(appName, env string) {
	printInfo("Start packing ...")
	Build(appName, env)

	bingoConfig, name := getBingoConfig(env)

	var zipName string
	if appName == "*" {
		zipName = "all-"
	} else {
		zipName = appName + "-"
	}
	zipName += time.Now().Format("060102")
	i := 1
	tmp := zipName + fmt.Sprintf("%02d.zip", i)
	for utils.IsFileExists(tmp) {
		i++
		tmp = zipName + fmt.Sprintf("%02d.zip", i)
	}
	zipName = tmp
	var ext string
	if runtime.GOOS == "windows" {
		ext = ".exe"
	} else {
		ext = ""
	}
	d, _ := os.Create(zipName)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	if appName == "*" {
		for _, app := range bingoConfig.Apps {
			if f, err := os.Open(filepath.Join(app.Package, app.Name+ext)); err == nil {
				if err = compressFile(f, w); err != nil {
					os.Remove(zipName)
					panic(err)
				}
				printSuccess("Compress: %s/%s Ok", app.Package, app.Name+ext)
			} else {
				os.Remove(zipName)
				panic(err)
			}
		}
	} else {
		app := bingoConfig.FindApp(appName)
		if f, err := os.Open(filepath.Join(app.Package, app.Name+ext)); err == nil {
			if err = compressFile(f, w); err != nil {
				os.Remove(zipName)
				panic(err)
			}
			printSuccess("Compress: %s/%s Ok", app.Package, app.Name+ext)
		} else {
			os.Remove(zipName)
			panic(err)
		}
	}
	if f, err := os.Open(name); err == nil {
		if err = compressFile(f, w); err != nil {
			os.Remove(zipName)
			panic(err)
		}
		printSuccess("Compress: %s Ok", name)
	} else {
		os.Remove(zipName)
		panic(err)
	}
	printSuccess("Publish done!")
}

func compressFile(file *os.File, w *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if !info.IsDir() {
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		writer, err := w.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
