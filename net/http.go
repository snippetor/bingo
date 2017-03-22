package net

import (
	"github.com/valyala/fasthttp"
)

func main() {
	c := fasthttp.Client{}
	c.Do()
}