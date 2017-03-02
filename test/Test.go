package main

import (
	"fmt"
)

func main() {
	//dir, _ := os.Getwd()
	//fmt.Printf(dir)
	//f, _ := os.OpenFile(dir+"/log.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	//defer f.Close()
	//l := log.New(f, "", log.Lmicroseconds)
	//l.Output(2, "[I] test")
	//l.Output(2, "[D] test1")

	c := make(chan string, 1000)
	go func() {
		for i := 0; i < 1000; i++ {
			c <- "goroutine1"
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			c <- "goroutine2"
		}
	}()

	go func() {
		for {
			select {
			case s := <- c:
				fmt.Println(s)
			}
		}
	}()

	for {
		if true {}
	}
}
