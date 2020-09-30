package main

import (

	// "time"

	// "io"

	// "fmt"
	"strings"

	"github.com/butuhanov/trading-helpers/server"
)


func main() {



 server.StartServer()



}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}
