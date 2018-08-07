package mvc

import "github.com/snippetor/bingo/route"

type Controller interface {
	Route(builder route.RouterBuilder)
}
