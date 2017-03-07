package main

import (
	"github.com/snippetor/bingo/log"
	"path/filepath"
	"os"
	"time"
	"fmt"
)

func main() {
	wd, _:= os.Getwd()
	logger := log.NewLogger(filepath.Join(wd, "log/example.ini"))

	go func() {
		t := time.NewTicker(3*time.Second)
		for {
			select {
			case <-t.C:
				fmt.Println("ticker1..")
				logger.I("test1 info========")
				logger.D("test1 debug========")
				logger.W("test1 warning========")
				logger.E("test1 err========")
			}
		}
	}()

	go func() {
		t := time.NewTicker(2*time.Second)
		for {
			select {
			case <-t.C:
				fmt.Println("ticker2..")
				logger.I("test2 info========")
				logger.D("test2 debug========")
				logger.W("test2 warning========")
				logger.E("test2 err========")
			}
		}
	}()

	c := make(chan bool, 1)
	<- c
}