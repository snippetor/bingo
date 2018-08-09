package mvc

import "github.com/snippetor/bingo/app"

type Controller interface {
	Route(builder app.RouterBuilder)
}
