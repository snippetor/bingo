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
	log.SetConfigFile(filepath.Join(wd, "test/log.config"))

	go func() {
		t := time.NewTicker(3*time.Second)
		for {
			select {
			case <-t.C:
				fmt.Println("ticker1..")
				log.I("test1 info========")
				log.D("test1 debug========")
				log.W("test1 warning========")
				log.E("test1 err========")
			}
		}
	}()

	go func() {
		t := time.NewTicker(2*time.Second)
		for {
			select {
			case <-t.C:
				fmt.Println("ticker2..")
				log.I("test2 info========")
				log.D("test2 debug========")
				log.W("test2 warning========")
				log.E("test2 err========")
			}
		}
	}()

	c := make(chan bool, 1)
	<- c
}