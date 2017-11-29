package node

import (
	"testing"
	"os"
	"fmt"
)

func TestNodeConfig(t *testing.T) {
	pwd, _ := os.Getwd()
	fmt.Println(pwd)
	Parse(pwd + "/../bingo.json")
	fmt.Println(config.Nodes[2].Config["port"].(float64))
}
