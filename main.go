package main

import (
	"github.com/butuhanov/trading-helpers/server"
	"github.com/pkg/profile"
)


func main() {

	defer profile.Start(profile.ProfilePath("/tmp")).Stop()


  server.StartServer()

}
