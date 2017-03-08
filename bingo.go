package bingo

import "github.com/snippetor/bingo/log"

var (
	fwLogger log.Logger
)
func init() {
	fwLogger = &log.NEW{}

	log.NewLogger()
}

func _frame_log()  {
	
}