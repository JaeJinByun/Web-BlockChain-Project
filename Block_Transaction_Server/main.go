package main

import (
	"Project"
)

func main() {

	go Project.StartTxServer()

	Project.StartBCServer()

}
